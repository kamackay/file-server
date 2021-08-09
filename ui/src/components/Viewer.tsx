import "promise-peek";
import React from "react";
import { isIOS, isSafari } from "react-device-detect";
import { RouteComponentProps, withRouter } from "react-router-dom";
import { BounceLoader as Loader } from "react-spinners";
import { WEBM_PATH } from "../constants";

interface Props extends RouteComponentProps<any> {}

interface State {
  path: string;
  width: number;
  height: number;
  downloading: boolean;
}

export default withRouter(
  class Browser extends React.Component<Props, State> {
    constructor(props: Props) {
      super(props);
      this.state = {
        path: this.props.location.pathname.replace(WEBM_PATH, ""),
        width: window.innerWidth * 0.75,
        height: window.innerHeight * 0.75,
        downloading: false,
      };
    }

    public componentWillUpdate(nextProps: Props, nextState: State) {
      if (this.props.location.pathname !== nextProps.location.pathname) {
        setTimeout(() => {
          this.setState(
            (prev) => ({
              ...prev,
              path: nextProps.location.pathname.replace(WEBM_PATH, ""),
            }),
            this.downloadIfIos
          );
        }, 1);
      }
    }

    public componentDidMount() {
      this.downloadIfIos();
    }

    public render() {
      const { path, width, height, downloading } = this.state;
      return (
        <div>
          {downloading ? (
            <Loader />
          ) : (
            <a
              style={{ display: "block" }}
              onClick={this.downloadFile}
              href="#"
            >
              Direct Link
            </a>
          )}
          {!isSafari && (
            <video
              autoPlay={true}
              loop={true}
              width={width}
              height={height}
              controls={true}
              playsInline={true}
            >
              <source src={path} type="video/mp4" />
              Sorry, your browser doesn't support embedded videos.{" "}
              <a href={path}>Download Here</a>
            </video>
          )}
        </div>
      );
    }

    private downloadIfIos = () => {
      if (isSafari || isIOS) {
        this.downloadFile();
      }
    };

    private downloadFile = () => {
      const { path } = this.state;
      this.setDownloading(true)();
      fetch(path)
        .then((response) => response.blob())
        .then((blob) => {
          const splitpath = path.split("/");
          const filename = splitpath[splitpath.length - 1];
          const url = window.URL.createObjectURL(blob);
          const a = document.createElement("a");
          a.href = url;
          a.download = filename;
          a.click();
        })
        .finally(this.setDownloading(false));
    };

    private setDownloading = (downloading: boolean) => () =>
      this.setState((prev) => ({ ...prev, downloading }));
  }
);
