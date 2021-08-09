import React from "react";
import ReactDOM from "react-dom";
import { Router, Route, Switch } from "react-router-dom";
import { createHashHistory } from "history";
import "./index.less";
import "antd/dist/antd.css";
import asyncComponent from "./components/asyncComponent";
import { ToastContainer } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";
import { Layout, PageHeader } from "antd";
import {BROWSE_PATH, WEBM_PATH} from "./constants";

const MainPage = asyncComponent(() => import("./components/MainPage"));
const Browser = asyncComponent(() => import("./components/Browser"));
const Viewer = asyncComponent(() => import("./components/Viewer"));

ReactDOM.render(
  <Router history={createHashHistory()}>
    <Layout style={{ paddingBottom: 10 }}>
      <ToastContainer position="top-right" />
      <Layout.Header>
        <PageHeader title="File Server" />
      </Layout.Header>
      <Layout.Content
        style={{ marginLeft: 10, marginRight: 10, padding: "0 50px" }}
      >
        <Switch>
          <Route path="/" exact={true} component={MainPage} />
          <Route path={`${BROWSE_PATH}**`} component={Browser} />
          <Route path={`${WEBM_PATH}**`} component={Viewer} />
        </Switch>
      </Layout.Content>
      <Layout.Footer style={{ textAlign: "center" }}>
        Keith MacKay ©2020
      </Layout.Footer>
    </Layout>
  </Router>,
  document.getElementById("root")
);
