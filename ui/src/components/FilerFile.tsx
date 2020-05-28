import { LockFilled as Lock } from "@ant-design/icons";
import { Card } from "antd";
import moment from "moment";
import React from "react";
import { Link, RouteComponentProps, withRouter } from "react-router-dom";
import urljoin from "url-join";

export default withRouter(
  (props: { file: FilerFile; path: string } & RouteComponentProps<any>) => {
    const { file, path } = props;
    const date = moment(file.lastUpdated / 1000000);
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
        <p>Type: {file.contentType}</p>
        <p>
          Last Updated{" "}
          {file.lastUpdated === 0 ? `Never` : date.format("Do MMM, YYYY")}
        </p>
      </Card>
    );
  }
);
