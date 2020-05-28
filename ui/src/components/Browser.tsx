import "promise-peek";

import { HomeOutlined as HomeIcon } from "@ant-design/icons";
import { Breadcrumb } from "antd";
import Axios from "axios";
import React from "react";
import { Link, RouteComponentProps, withRouter } from "react-router-dom";
import { BarLoader as Loader } from "react-spinners";
import { toast } from "react-toastify";
import urljoin from "url-join";

import { BROWSE_PATH } from "../constants";
import FilerFile from "./FilerFile";

type Contents = FilerFile[] | string;
type Type = "folder" | "file";

interface Props extends RouteComponentProps<any> {}

interface State {
  pathname?: string;
  type: Type;
  path?: string[];
  contents?: Contents;
}

export default withRouter(
  class Browser extends React.Component<Props, State> {
    constructor(props: Props) {
      super(props);
      this.state = {
        type: "folder",
      };
    }

    public componentDidMount() {
      this.loadData();
    }

    public componentWillUpdate(nextProps: Props, nextState: State) {
      if (this.props.location.pathname !== nextProps.location.pathname) {
        setTimeout(() => {
          this.setState(
            (prev) => ({ ...prev, contents: undefined }),
            this.loadData
          );
        }, 1);
      }
    }

    public render() {
      const { contents } = this.state;
      return !!contents ? <this.renderContents /> : <Loader />;
    }

    private loadData = () => {
      const path = this.props.location.pathname.replace(BROWSE_PATH, "");
      const split = path.split(`/`);
      let type: Type = "folder";
      Axios.get(path, {
        headers: { "Get-Folder": `true` },
      })
        .peek((r) => {
          switch (r.headers["type"]) {
            case "file":
              type = "file";
            default:
              break;
          }
        })
        .then((r) => r.data)
        .then((data: Contents) => {
          this.setState((prev) => ({
            ...prev,
            type,
            contents: data,
            pathname: path,
            path: split.filter((o) => !!o),
          }));
        })
        .catch((err) => {
          console.warn(err);
          toast(`Error Getting Path`);
        });
    };

    private renderContents = () => {
      switch (this.state.type) {
        case "file":
          return <this.renderFile />;
        default:
        case "folder":
          return <this.renderFolder />;
      }
    };

    private renderFile = () => {
      const contents = this.state.contents as string;
      return <div>{contents.length}</div>;
    };

    private renderFolder = () => {
      const contents = this.state.contents as FilerFile[];
      const { path, pathname } = this.state;
      return (
        <div>
          <Breadcrumb>
            {[
              <Breadcrumb.Item key={`home`}>
                <Link to="/">
                  <HomeIcon />
                </Link>
              </Breadcrumb.Item>,
              ...(path || []).map((part, x) => {
                const route = [
                  BROWSE_PATH,
                  ...(path || []).slice(0, x + 1),
                ].join(`/`);
                return (
                  <Breadcrumb.Item key={`part-${part}-${x}`}>
                    <Link to={route}>{part}</Link>
                  </Breadcrumb.Item>
                );
              }),
            ]}

            <div style={{ display: "block" }}>
              {contents!.map((file, x) => {
                return (
                  <FilerFile
                    key={`${x}`}
                    path={urljoin(BROWSE_PATH, pathname || "")}
                    file={file}
                  />
                );
              })}
            </div>
          </Breadcrumb>
        </div>
      );
    };
  }
);
