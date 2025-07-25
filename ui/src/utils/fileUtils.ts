/**
 * Utility functions for working with files
 */

/**
 * Fetches a text file from the server
 * @param filePath - The path to the file relative to the public directory
 * @returns A promise that resolves to the text content of the file
 */
export const fetchTextFile = async (filePath: string): Promise<string> => {
  try {
    const response = await fetch(filePath);
    
    if (!response.ok) {
      throw new Error(`Failed to fetch file: ${filePath}`);
    }
    
    return await response.text();
  } catch (error) {
    console.error('Error fetching text file:', error);
    throw error;
  }
};