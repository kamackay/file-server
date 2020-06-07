declare interface ConversionJob {
  id: string;
  inputFile: string;
  outputFile: string;
  status: number;
  progress: number;
  error: any;
}
