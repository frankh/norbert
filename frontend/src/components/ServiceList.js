import React from 'react';
import { Query, graphql } from "react-apollo";
import './App.css';

import { GET_SERVICES } from "../queries";


class ServiceListItem extends React.PureComponent {
    render() {
        return (
            <div className="serviceListItem">
                {this.props.name}
            </div>
        )
    }
}


class ServiceList extends React.PureComponent {
    render() {
        return (
            <Query query={GET_SERVICES}>
                {({ loading, error, data, subscribeToMore }) => {
                    if (error) {
                        return (
                            <div>Error</div>
                        )
                    }
                    if (loading) {
                        return (
                            <div>Loading</div>
                        )
                    }
                    return (
                        <div className="serviceList">
                            {data.services.map((service) => (
                                <ServiceListItem name={service.name} />
                            ))}
                        </div>
                    )
                }}
            </Query>
        )
    }
}

export default graphql(GET_SERVICES, { name: "getServices" })(
    ServiceList
);
