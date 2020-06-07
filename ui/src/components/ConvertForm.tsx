import { UploadOutlined } from "@ant-design/icons";
import { Button, Col, Form, Input, Progress, Row } from "antd";
import axios from "axios";
import React from "react";
import { toast } from "react-toastify";
import { makeAuthHeader } from "../utils/utils";

interface Props {
  creds?: Creds;
}

interface State {
  url: string;
  uploadUrl: string;
  loading: boolean;
  job?: ConversionJob;
}

export default class ConvertForm extends React.Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = {
      url: "/",
      uploadUrl: "/",
      loading: false,
    };
  }

  public componentWillUnmount() {
    // No-op
  }

  public render() {
    const { url, uploadUrl, loading, job } = this.state;
    const { creds } = this.props;
    return (
      <Form
        onFinish={(values: { url: string; filePath: string }) => {
          this.setState((prev) => ({ ...prev, loading: true }));
          axios
            .post(
              values.filePath,
              { output: values.url },
              {
                headers: {
                  "Content-Type": `file/convert`,
                  ...makeAuthHeader(creds),
                },
              }
            )
            .then((r) => r.data)
            .then((data: { job: string }) => {
              const { job } = data;
              this.checkOnJob(job);
            })
            .catch((err) => {
              console.warn(err);
              toast(`Could not Convert URL`, { type: "error" });
              this.setState((prev) => ({ ...prev, loading: false }));
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
                    const url = e.target.value;
                    this.setState((prev) => ({ ...prev, url }));
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
              const uploadUrl = e.target.value;
              this.setState((prev) => ({ ...prev, uploadUrl }));
            }}
          />
        </Form.Item>

        <Form.Item>
          {!!job && (
            <div style={{ margin: 5, padding: 5 }}>
              <Progress
                strokeColor={{
                  from: "#108ee9",
                  to: "#87d068",
                }}
                percent={job.progress}
                status="active"
              />
            </div>
          )}
          <Button loading={loading} type="primary" htmlType="submit">
            <UploadOutlined /> Convert
          </Button>
        </Form.Item>
      </Form>
    );
  }

  private checkOnJob = (job: string) => {
    axios
      .post(`/`, `${job}`, {
        headers: {
          "Content-Type": `job/progress`,
          ...makeAuthHeader(this.props.creds),
        },
      })
      .then((r) => r.data)
      .then((data: ConversionJob) => {
        this.setState((prev) => ({ ...prev, job: data }));
        if (data.status !== 0 && this.state.loading) {
          this.setState((prev) => ({
            ...prev,
            loading: false,
            job: undefined,
          }));
          if (data.status === 1) {
            toast(`Could not Convert URL`, { type: "error" });
          } else if (data.status === 2) {
            toast(`Successfully Converted`, { type: "success" });
          }
        } else {
          setTimeout(() => this.checkOnJob(job), 100);
        }
      })
      .catch((err) => {
        console.warn(err);
      });
  };
}
