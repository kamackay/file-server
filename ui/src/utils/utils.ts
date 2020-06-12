export const makeAuthHeader = (creds?: Creds) =>
  !!creds
    ? {
        Authorization: `Basic ${btoa(`${creds.username}:${creds.password}`)}`,
      }
    : { Authorization: `` };

const unitSize = 1024;

export const humanizeBytes = (bytes: number): string => {
  let unit = 0;
  let bs = bytes;
  while (bs > unitSize) {
    bs /= unitSize;
    unit++;
  }

  return `${formatDecimal(bs, 3)}${
    [" bytes", "kB", "MB", "GB", "TB", "PB"][unit]
  }`;
};

export const formatDecimal = (n: number, precision: number): string => {
  const factor = Math.pow(10, precision);
  return `${Math.round(n * factor) / factor}`;
};
