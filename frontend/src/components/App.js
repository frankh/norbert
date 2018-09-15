import React, { Component } from 'react';
import { graphql } from "react-apollo";
import './App.css';
import ServiceList from './ServiceList.js';

import { GET_SERVICES } from "../queries";


class App extends Component {
  render() {
    return (
      <div className="App">
          <ServiceList />
      </div>
    );
  }
}

export default graphql(GET_SERVICES, { name: "getServices" })(
    App
);
