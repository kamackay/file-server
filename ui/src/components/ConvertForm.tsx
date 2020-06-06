import { UploadOutlined } from "@ant-design/icons";
import { Button, Col, Form, Input, Row } from "antd";
import axios from "axios";
import React from "react";
import { toast } from "react-toastify";
import { makeAuthHeader } from "../utils/utils";

export default (props: { creds?: Creds }) => {
  const [url, setUrl] = React.useState<string>(`/file.txt`);
  const [uploadUrl, setUploadUrl] = React.useState<string>(``);
  const [loading, setLoading] = React.useState<boolean>(false);

  return (
    <Form
      onFinish={(values: { url: string; filePath: string }) => {
        setLoading(true);
        axios
          .post(
            values.filePath,
            { output: values.url },
            {
              headers: {
                "Content-Type": `file/convert`,
                ...makeAuthHeader(props.creds),
              },
            }
          )
          .then((r) => r.data)
          .then((data) => {
            toast(`Successfully Uploaded!`, { type: "success" });
          })
          .catch((err) => {
            console.warn(err);
            toast(`Could not upload URL`, { type: "error" });
          })
          .finally(() => {
            setLoading(false);
          });
      }}
    >
      <h3>Conversion</h3>
      <h5>
        The server will attempt to use FFMPEG to convert from one file to the
        other
      </h5>

      <Row gutter={24}>
        <Col span={6} />
        <Col
          span={12}
          children={
            <Form.Item
              label="Input Path"
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
        label="Output Path"
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
        <Button loading={loading} type="primary" htmlType="submit">
          <UploadOutlined /> Convert
        </Button>
      </Form.Item>
    </Form>
  );
};
