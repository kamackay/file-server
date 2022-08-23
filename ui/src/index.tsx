import { Layout, PageHeader } from "antd";
import "antd/dist/antd.css";
import { createHashHistory } from "history";
import React from "react";
import ReactDOM from "react-dom";
import { Route, Router, Switch } from "react-router-dom";
import { ToastContainer } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";
import asyncComponent from "./components/asyncComponent";
import "./index.less";

const MainPage = asyncComponent(() => import("./components/MainPage"));
const Browser = asyncComponent(() => import("./components/Browser"));

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
          <Route path="/browse**" component={Browser} />
        </Switch>
      </Layout.Content>
      <Layout.Footer style={{ textAlign: "center" }}>
        Keith MacKay ©{new Date().getFullYear()}
      </Layout.Footer>
    </Layout>
  </Router>,
  document.getElementById("root")
);
