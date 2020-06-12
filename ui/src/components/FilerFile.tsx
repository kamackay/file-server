import { LockFilled as Lock } from "@ant-design/icons";
import { Card } from "antd";
import moment from "moment";
import React from "react";
import { Link, RouteComponentProps, withRouter } from "react-router-dom";
import urljoin from "url-join";
import { humanizeBytes } from "../utils/utils";

export default withRouter(
  (props: { file: FilerFile; path: string } & RouteComponentProps<any>) => {
    const { file, path } = props;
    return (
      <Card
        size="small"
        title={
          <p>
            {file.name} {file.protected && <Lock />}
          </p>
        }
        style={{ display: "inline-block", width: 300, margin: 10 }}
        extra={<Link to={`${urljoin(path, file.name)}`} children="View" />}
      >
        <RenderInners file={file} />
      </Card>
    );
  }
);

const RenderFolderInners = (props: { file: FilerFile }): JSX.Element => {
  const { file } = props;
  return (
    <>
      <p>Type: Folder</p>
      {file.count >= 0 && <p>Items: {file.count}</p>}
      {file.size >= 0 && <p>Size: {humanizeBytes(file.size)}</p>}
    </>
  );
};

const RenderFileInners = (props: { file: FilerFile }): JSX.Element => {
  const { file } = props;
  const date = moment(file.lastUpdated / 1000000);
  return (
    <>
      <p>Type: {file.contentType}</p>
      <p>
        Last Updated{" "}
        {file.lastUpdated === 0 ? `Never` : date.format("Do MMM, YYYY")}
      </p>
      {file.size >= 0 && <p>Size: {humanizeBytes(file.size)}</p>}
    </>
  );
};

const RenderInners = (props: { file: FilerFile }): JSX.Element => {
  const { file } = props;
  switch (props.file.contentType) {
    case "folder":
      return <RenderFolderInners file={file} />;
    default:
      return <RenderFileInners file={file} />;
  }
};
