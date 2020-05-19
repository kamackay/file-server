import React from "react";
import { UploadOutlined } from "@ant-design/icons";
import { Row, Form, Col, Input, Upload } from "antd";
import { toast } from "react-toastify";
import { makeAuthHeader } from "../utils/utils";

export default (props: { creds?: Creds }) => {
  const [url, setUrl] = React.useState<string>(`/file.txt`);
  return (
    <Form>
      <h3>Upload a File</h3>

      <Row gutter={24}>
        <Col span={6} />
        <Col
          span={12}
          children={
            <Form.Item
              label="Path"
              name="path"
              initialValue={url}
              rules={[{ required: true, message: "Please input A Path!" }]}
            >
              <Input
                value={url}
                onChange={(e) => {
                  const url = e.target.value;
                  setUrl(url);
                }}
              />
            </Form.Item>
          }
        />
        <Col span={6} />
      </Row>

      <Upload.Dragger
        method="PUT"
        action={url}
        onChange={(state) => {
          if (state.file.status === "done") {
            toast(`Uploaded Successfully!`, { type: "success" });
          }
        }}
        headers={{
          ...makeAuthHeader(props.creds),
        }}
      >
        <p className="ant-upload-drag-icon">
          <UploadOutlined /> Click or Drag To Upload
        </p>
      </Upload.Dragger>
    </Form>
  );
};
