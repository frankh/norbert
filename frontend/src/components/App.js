import React, { Component } from 'react';
import { graphql } from "react-apollo";
import { Layout } from "antd";
import { Menu, Icon, Button } from 'antd';
import ServiceList from './ServiceList.js';
import { Route, Switch, withRouter } from "react-router-dom";
import CheckPage from './CheckPage.js';

import { GET_SERVICES } from "../queries";

import "antd/dist/antd.css";
import './App.css';

const { Header, Sider, Content } = Layout;


class App extends Component {
  render() {
    return (
      <div className="App">
        <Layout>
          <Sider>
            <div class="logo-title">
              <img src="logo.png"></img>
              <h1>Norbert</h1>
            </div>
            <Menu
              defaultSelectedKeys={['services']}
              mode="inline"
              theme="light"
            >
              <Menu.Item key="services">
                <Icon type="deployment-unit" theme="outlined" />
                <span>Services</span>
              </Menu.Item>
              <Menu.Item key="checks">
                <Icon type="fork" theme="outlined" />
                <span>Checks</span>
              </Menu.Item>
            </Menu>
          </Sider>
          <Layout style={{backgroundColor: "inherit"}}>
              <Header />
              <Content>
                  <Switch>
                      <Route path="/check/:id/" component={CheckPage} />
                      <Route path="/" component={ServiceList}/>
                  </Switch>
              </Content>
          </Layout>
        </Layout>
      </div>
    );
  }
}

export default graphql(GET_SERVICES, { name: "getServices" })(
    withRouter(App)
);
