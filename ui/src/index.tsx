import React from "react";
import ReactDOM from "react-dom";
import { HashRouter as Router, Route } from "react-router-dom";
import "./index.less";
import "antd/dist/antd.css";
import asyncComponent from "./components/asyncComponent";
import { Layout } from "antd";

const MainPage = asyncComponent(() => import("./components/MainPage"));

ReactDOM.render(
  <Router>
    <Layout style={{ paddingBottom: 10 }}>
      <Layout.Content style={{ marginLeft: 10, marginRight: 10 }}>
        <Route path="/" exact={true} component={MainPage} />
      </Layout.Content>
    </Layout>
  </Router>,
  document.getElementById("root")
);
