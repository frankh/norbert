import React from 'react';
import { Query, graphql } from "react-apollo";
import { List } from "antd";

import { GET_SERVICES, SERVICE_SUBSCRIPTION } from "../queries";


const StatusToOrder = {
    "Ok": 0,
    "Failed": 1,
}

const OrderToStatus = {
    NaN: "Ok",
    0: "Ok",
    1: "Failed",
}

class CheckSummary extends React.PureComponent {
    render() {
        return (
            <div className={"checkSummary " + this.props.check.status} title={this.props.check.name}>
            </div>
        )
    }
}

class ServiceListItem extends React.PureComponent {
    componentDidMount() {
        this.onUpdate = this.props.subscribe({
            document: SERVICE_SUBSCRIPTION,
            variables: { serviceName: this.props.service.name },
            updateQuery: (prev, { subscriptionData: { data } }) => {
                const newService = data.serviceChanged
                for( var i = 0; i < prev.services.length; i++ ) {
                    if( prev.services[i].name === newService.name ) {
                        prev.services[i] = newService
                    }
                }
                return prev
            }
        })
    }

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
                    {checks.map(check => <CheckSummary check={check} />)}
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
                                    <ServiceListItem key={service.name} service={service} subscribe={subscribeToMore}/>
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
