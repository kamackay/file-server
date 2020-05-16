import React from "react";
import ReactDOM from "react-dom";
import { HashRouter as Router, Route } from "react-router-dom";
import "./index.less";
import "antd/dist/antd.css";
import asyncComponent from "./components/asyncComponent";
import { ToastContainer } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";
import { Layout } from "antd";

const MainPage = asyncComponent(() => import("./components/MainPage"));

ReactDOM.render(
  <Router>
    <Layout style={{ paddingBottom: 10 }}>
      <ToastContainer position="top-right" />
      <Layout.Content
        style={{ marginLeft: 10, marginRight: 10, padding: "0 50px" }}
      >
        <Route path="/" exact={true} component={MainPage} />
      </Layout.Content>
      <Layout.Footer style={{ textAlign: "center" }}>
        Keith MacKay ©2020
      </Layout.Footer>
    </Layout>
  </Router>,
  document.getElementById("root")
);
