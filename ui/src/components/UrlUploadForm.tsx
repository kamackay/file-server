import { UploadOutlined } from "@ant-design/icons";
import { Button, Col, Form, Input, Row } from "antd";
import axios from "axios";
import React from "react";
import { toast } from "react-toastify";
import { makeAuthHeader } from "../utils/utils";

export default (props: { creds?: Creds }) => {
  const [url, setUrl] = React.useState<string>(`/file.txt`);
  const [uploadUrl, setUploadUrl] = React.useState<string>(``);

  return (
    <Form
      onFinish={(values: { url: string; filePath: string }) => {
        axios
          .post(values.filePath, values.url, {
            headers: {
              "Content-Type": `text/plain`,
              ...makeAuthHeader(props.creds),
            },
          })
          .then((r) => r.data)
          .then((data) => {
            toast(`Successfully Uploaded!`, { type: "success" });
          })
          .catch((err) => {
            console.warn(err);
            toast(`Could not upload URL`, { type: "error" });
          });
      }}
    >
      <h3>Upload a Remote File</h3>
      <h5>
        The Server will download the proxided URL to its filesystem and it will
        be available as a regular file
      </h5>

      <Row gutter={24}>
        <Col span={6} />
        <Col
          span={12}
          children={
            <Form.Item
              label="File Path"
              name="filePath"
              initialValue={url}
              rules={[{ required: true, message: "Please Input A Path!" }]}
            >
              <Input
                value={url}
                onChange={(e) => {
                  setUrl(e.target.value);
                }}
              />
            </Form.Item>
          }
        />
        <Col span={6} />
      </Row>

      <Form.Item
        label="URL"
        name="url"
        initialValue={uploadUrl}
        rules={[{ required: true, message: "Please input A Path!" }]}
      >
        <Input
          value={uploadUrl}
          onChange={(e) => {
            setUploadUrl(e.target.value);
          }}
        />
      </Form.Item>

      <Form.Item>
        <Button type="primary" htmlType="submit">
          <UploadOutlined /> Upload
        </Button>
      </Form.Item>
    </Form>
  );
};
