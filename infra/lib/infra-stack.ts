import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as s3 from 'aws-cdk-lib/aws-s3';
import * as s3deploy from 'aws-cdk-lib/aws-s3-deployment';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as apigwv2 from 'aws-cdk-lib/aws-apigatewayv2';
import { HttpLambdaIntegration } from 'aws-cdk-lib/aws-apigatewayv2-integrations';
import * as s3n from 'aws-cdk-lib/aws-s3-notifications';
import * as iam from 'aws-cdk-lib/aws-iam';

export class InfraStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // 1. S3 バケットの作成

    const removalPolicy = cdk.RemovalPolicy.DESTROY;
    const autoDeleteObjects = true;

    // Inputバケット（画像アップロード用：1日で自動削除）
    const inputBucket = new s3.Bucket(this, 'InputBucket', {
      removalPolicy,
      autoDeleteObjects,
      lifecycleRules: [{ expiration: cdk.Duration.days(1) }],
      cors: [{
        allowedMethods: [s3.HttpMethods.PUT],
        allowedOrigins: ['*'],
        allowedHeaders: ['*'],
      }],
    });

    // Outputバケット（生成したExcel保存用：1日で自動削除）
    const outputBucket = new s3.Bucket(this, 'OutputBucket', {
      removalPolicy,
      autoDeleteObjects,
      lifecycleRules: [{ expiration: cdk.Duration.days(1) }],
    });

    // Websiteバケット（フロントエンドの静的ファイルホスティング用）
    const websiteBucket = new s3.Bucket(this, 'WebsiteBucket', {
      removalPolicy,
      autoDeleteObjects,
      websiteIndexDocument: 'index.html',
      publicReadAccess: true,
      blockPublicAccess: new s3.BlockPublicAccess({
        blockPublicAcls: false,
        blockPublicPolicy: false,
        ignorePublicAcls: false,
        restrictPublicBuckets: false,
      }),
    });

    // 差し替え用 Website HTML を Homepageバケットに挿入
    new s3deploy.BucketDeployment(this, 'DeployWebsite', {
      sources: [s3deploy.Source.asset('./test-assets')],
      destinationBucket: websiteBucket,
    });

    // 2. Lambda 関数の作成（Goランタイム）

    // PresignHandler（署名付きURL発行）
    const presignHandler = new lambda.Function(this, 'PresignHandlerv2', {
      runtime: lambda.Runtime.PROVIDED_AL2023,
      handler: 'bootstrap',
      architecture: lambda.Architecture.ARM_64,
      code: lambda.Code.fromAsset('../backend/bin/presign'),
      environment: {
        INPUT_BUCKET_NAME: inputBucket.bucketName,
      },
    });

    // MainHandler（Textract解析・JSON生成・Slack通知）
    const mainHandler = new lambda.Function(this, 'MainHandlerv2', {
      runtime: lambda.Runtime.PROVIDED_AL2023,
      handler: 'bootstrap',
      architecture: lambda.Architecture.ARM_64,
      code: lambda.Code.fromAsset('./test-assets/dummy-lambda'),
      timeout: cdk.Duration.seconds(30),
      environment: {
        OUTPUT_BUCKET_NAME: outputBucket.bucketName,
        OUTPUT_FORMAT: 'json', 
        SLACK_WEBHOOK_URL: process.env.SLACK_WEBHOOK_URL || '',
      },
    });

    // 3. 権限（IAM）と トリガー（Event）の設定

    // PresignHandlerには、Inputバケットへ「書き込む」権限のみ付与
    inputBucket.grantWrite(presignHandler);

    // MainHandlerには、Inputから「読み込む」権限、Outputへ「書き込む」権限を付与
    inputBucket.grantRead(mainHandler);
    outputBucket.grantWrite(mainHandler);

    // MainHandlerにTextractの実行権限を付与
    mainHandler.addToRolePolicy(new iam.PolicyStatement({
      actions: ['textract:AnalyzeDocument', 'textract:DetectDocumentText'],
      resources: ['*'],
    }));

    // Inputバケットに画像が入ったら MainHandler を自動起動
    inputBucket.addEventNotification(
      s3.EventType.OBJECT_CREATED,
      new s3n.LambdaDestination(mainHandler)
    );

    // 4. API Gateway の構築 (HTTP API)

    const api = new apigwv2.HttpApi(this, 'JpegToJsonHttpApi', {
      apiName: 'Jpeg To Json HTTP API',
      corsPreflight: {
        allowOrigins: ['*'],
        allowMethods: [apigwv2.CorsHttpMethod.ANY],
      },
    });

    // GET /presign で PresignHandler を呼び出す
    const presignIntegration = new HttpLambdaIntegration('PresignIntegration', presignHandler);
    api.addRoutes({
      path: '/presign',
      methods: [apigwv2.HttpMethod.GET],
      integration: presignIntegration,
    });

    // 5. デプロイ後の出力も更新
    new cdk.CfnOutput(this, 'ApiEndpoint', { value: api.apiEndpoint });
  }
}
