import React from 'react';
import { Query, graphql } from "react-apollo";
import { List, Icon } from "antd";
import './App.css';

import { GET_SERVICES } from "../queries";


class ServiceListItem extends React.PureComponent {
    render() {
        const service = this.props.service;
        return (
            <div className="serviceListItem">
                <List.Item>
                    <List.Item.Meta
                      title={service.name}
                      description={<a href={service.url}>{service.url}</a>}

                    />
                    <Icon type="undo" theme="outlined" spin className="reverse"/>
                </List.Item>
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
                            <List
                                dataSource={data.services}
                                renderItem={service => (
                                    <ServiceListItem key={service.name} service={service} />
                                )}
                            />
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
