import React from "react";
import { Form, Input, Button, Upload } from "antd";
import Cookies from "js-cookie";
import { UploadOutlined } from "@ant-design/icons";
import { withRouter, RouteComponentProps } from "react-router-dom";

interface Props extends RouteComponentProps<any> {}

interface State {
  creds?: Creds;
  url: string;
}

interface Creds {
  username: string;
  password: string;
}

export default withRouter(
  class MainPage extends React.Component<Props, State> {
    constructor(props: Props) {
      super(props);
      this.state = {
        url: `/file`,
        creds: Cookies.getJSON(`creds`),
      };
    }

    public render() {
      return (
        <div>
          <h1>Hey!</h1>
          <h3>Use this page to upload files to the Static File Server</h3>

          <Form
            wrapperCol={{ span: 16 }}
            onFinish={(values: Creds) => {
              Cookies.set(`creds`, values);
              this.setState((prev) => ({ ...prev, creds: values }));
            }}
            initialValues={this.state.creds}
            style={{ textAlign: "center" }}
          >
            <Form.Item
              label="Username"
              name="username"
              rules={[
                { required: true, message: "Please input your username!" },
              ]}
            >
              <Input />
            </Form.Item>

            <Form.Item
              label="Password"
              name="password"
              rules={[
                { required: true, message: "Please input your password!" },
              ]}
            >
              <Input.Password />
            </Form.Item>

            <Form.Item>
              <Button type="primary" htmlType="submit">
                Save
              </Button>
            </Form.Item>
          </Form>

          <Form wrapperCol={{ span: 16, style: { textAlign: "center" } }}>
            <h3>Upload a File</h3>

            <Form.Item
              label="Path"
              name="path"
              initialValue={this.state.url}
              rules={[{ required: true, message: "Please input A Path!" }]}
            >
              <Input
                value={this.state.url}
                onChange={(e) => {
                  const url = e.target.value;
                  console.log(url);
                  this.setState((prev) => ({ ...prev, url }));
                }}
              />
            </Form.Item>

            <Upload
              method="PUT"
              action={this.state.url}
              headers={{
                Authorization: `Basic ${btoa(
                  `${this.state.creds?.username}:${this.state.creds?.password}`
                )}`,
              }}
            >
              <Button>
                <UploadOutlined /> Click to Upload
              </Button>
            </Upload>
          </Form>
        </div>
      );
    }
  }
);
