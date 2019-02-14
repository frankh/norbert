import React from 'react';
import { Query, graphql } from "react-apollo";
import { List, Icon } from "antd";

import { GET_SERVICES } from "../queries";


const StatusToOrder = {
    "Ok": 0,
    "Failed": 1,
}

const OrderToStatus = {
    NaN: "Ok",
    0: "Ok",
    1: "Failed",
}

class ServiceListItem extends React.PureComponent {
    render() {
        const service = this.props.service;
        const checks = service.checks;
        const status = OrderToStatus[Math.max(...checks.map(check => StatusToOrder[check.status]))];

        return (
            <div className={`serviceListItem state${status}`}>
                <List.Item>
                    <List.Item.Meta
                      title={service.name}
                      description={<a href={service.url}>{service.url}</a>}
                    />
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
