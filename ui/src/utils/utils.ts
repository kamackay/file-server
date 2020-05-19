export const makeAuthHeader = (creds?: Creds) =>
  !!creds
    ? {
        Authorization: `Basic ${btoa(`${creds.username}:${creds.password}`)}`,
      }
    : { Authorization: `` };
