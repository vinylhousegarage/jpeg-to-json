export const compressImage = async (file: File): Promise<Blob> => {
  console.log("Starting image compression:", file.name);
  return file; 
};
