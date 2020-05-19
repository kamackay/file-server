import React from "react";
import { Form, Input, Button, Row, Col, Tabs } from "antd";
import Cookies from "js-cookie";
import { withRouter, RouteComponentProps } from "react-router-dom";
import FileUploadForm from "./FileUploadForm";
import UrlUploadForm from "./UrlUploadForm";
import CreateProxyForm from "./CreateProxyForm";

interface Props extends RouteComponentProps<any> {}

interface State {
  creds?: Creds;
}

export default withRouter(
  class MainPage extends React.Component<Props, State> {
    constructor(props: Props) {
      super(props);
      this.state = {
        creds: Cookies.getJSON(`creds`),
      };
    }

    public render() {
      return (
        <div>
          <h1>Hey!</h1>
          <h3>Use this page to upload files to the Static File Server</h3>

          <Form
            onFinish={(values: Creds) => {
              Cookies.set(`creds`, values);
              this.setState((prev) => ({ ...prev, creds: values }));
            }}
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
                    children={<Input />}
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
                    children={<Input.Password />}
                  />
                }
              />
            </Row>

            <Form.Item>
              <Button type="primary" htmlType="submit">
                Save
              </Button>
            </Form.Item>
          </Form>

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
          </Tabs>
        </div>
      );
    }
  }
);
