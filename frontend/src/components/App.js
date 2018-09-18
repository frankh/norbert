import React, { Component } from 'react';
import { graphql } from "react-apollo";
import { Layout } from "antd";
import ServiceList from './ServiceList.js';
import { Route, Switch, withRouter } from "react-router-dom";
import CheckPage from './CheckPage.js';

import { GET_SERVICES } from "../queries";

import "antd/dist/antd.css";
import './App.css';

const { Header, Footer, Content } = Layout;


class App extends Component {
  render() {
    return (
      <div className="App">
            <Layout style={{backgroundColor: "inherit"}}>
                <Header />
                <Content>
                    <Switch>
                        <Route path="/check/:id/" component={CheckPage} />
                        <Route path="/" component={ServiceList}/>
                    </Switch>
                </Content>
                <Footer />
            </Layout>
      </div>
    );
  }
}

export default graphql(GET_SERVICES, { name: "getServices" })(
    withRouter(App)
);
