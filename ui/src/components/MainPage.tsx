import { FolderViewOutlined as BrowseIcon } from "@ant-design/icons";
import { Button, Col, Form, Input, Row, Tabs } from "antd";
import Cookies from "js-cookie";
import React from "react";
import { RouteComponentProps, withRouter } from "react-router-dom";
import ConvertForm from "./ConvertForm";
import CreateProxyForm from "./CreateProxyForm";
import FileUploadForm from "./FileUploadForm";
import UrlUploadForm from "./UrlUploadForm";

interface Props extends RouteComponentProps<any> {}

interface State {
  creds: Creds;
}

export default withRouter(
  class MainPage extends React.Component<Props, State> {
    constructor(props: Props) {
      super(props);
      this.state = {
        creds: Cookies.getJSON(`creds`) || { username: "", password: "" },
      };
    }

    private setCred = (key: keyof Creds) => (
      event: React.ChangeEvent<HTMLInputElement>
    ) => {
      const value = event.target.value;
      this.setState(
        (prev) => ({
          ...prev,
          creds: { ...prev.creds, [key]: value },
        }),
        () => {
          Cookies.set(`creds`, this.state.creds);
        }
      );
    };

    public render() {
      return (
        <div>
          <h1>Hey!</h1>
          <h3>Use this page to upload files to the Static File Server</h3>

          <Form
            initialValues={this.state.creds}
            style={{ textAlign: "center" }}
          >
            <Row gutter={24}>
              <Col
                span={12}
                children={
                  <Form.Item
                    label="Username"
                    name="username"
                    rules={[
                      {
                        required: true,
                        message: "Please input your username!",
                      },
                    ]}
                    children={<Input onChange={this.setCred("username")} />}
                  />
                }
              />

              <Col
                span={12}
                children={
                  <Form.Item
                    label="Password"
                    name="password"
                    rules={[
                      {
                        required: true,
                        message: "Please input your password!",
                      },
                    ]}
                    children={
                      <Input.Password onChange={this.setCred("password")} />
                    }
                  />
                }
              />
            </Row>
          </Form>

          <Button
            icon={<BrowseIcon />}
            onClick={() => this.props.history.push(`/browse/`)}
            children="Browse"
          />

          <Tabs defaultActiveKey="fileUpload">
            <Tabs.TabPane
              tab="File Upload"
              key="fileUpload"
              children={<FileUploadForm creds={this.state.creds} />}
            />
            <Tabs.TabPane
              tab="URL Upload"
              key="urlUpload"
              children={<UrlUploadForm creds={this.state.creds} />}
            />
            <Tabs.TabPane
              tab="Proxy"
              key="proxy"
              children={<CreateProxyForm creds={this.state.creds} />}
            />
            <Tabs.TabPane
              tab="Convert"
              key="convert"
              children={<ConvertForm creds={this.state.creds} />}
            />
          </Tabs>
        </div>
      );
    }
  }
);
