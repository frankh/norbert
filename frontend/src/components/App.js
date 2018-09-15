import React, { Component } from 'react';
import { graphql } from "react-apollo";
import { Layout } from "antd";
import ServiceList from './ServiceList.js';

import { GET_SERVICES } from "../queries";

import './App.css';
import "antd/dist/antd.css";

const { Header, Footer, Content } = Layout;


class App extends Component {
  render() {
    return (
      <div className="App">
            <Layout>
                <Header />
                <Content>
                    <ServiceList />
                </Content>
                <Footer />
            </Layout>
      </div>
    );
  }
}

export default graphql(GET_SERVICES, { name: "getServices" })(
    App
);
