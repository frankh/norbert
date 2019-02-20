import React from 'react';
import { Query, graphql } from "react-apollo";
import { List } from "antd";
import shallowEqual from 'shallowequal'

/** @jsx jsx */
import { jsx, css, keyframes } from '@emotion/core'
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
    constructor(props) {
        super(props);

        this.state = {
            animRotate: keyframes`
                0% {transform: rotate(0deg);}
                100% {transform: rotate(360deg);}
            `,
            animOpacity: keyframes`
                0% {opacity: 1;}
                50%, 100% {opacity: 0;}
            `
        }
    }

    componentDidMount() {
        this.setState({
            animRotate: keyframes`
                0% {transform: rotate(0deg);}
                100% {transform: rotate(360deg);} //${Math.random()}
            `,
            animOpacity: keyframes`
                0% {opacity: 1;}
                50%, 100% {opacity: 0;} //${Math.random()}
            `
        })
    }

    componentWillReceiveProps(nextProps) {
        if (!shallowEqual(this.props.check, nextProps.check)) {
            this.setState({
                animRotate: keyframes`
                    0% {transform: rotate(0deg);}
                    100% {transform: rotate(360deg);} //${Math.random()}
                `,
                animOpacity: keyframes`
                    0% {opacity: 1;}
                    50%, 100% {opacity: 0;} //${Math.random()}
                `
            })
        }
    }

    render() {
        const colorMap = {
            'Ok': '#85e102',
            'Failed': '#e2312c',
        }
        const borderColorMap = {
            'Ok': '#b3ef5a',
            'Failed': '#ef807b',
        }
        const color = colorMap[this.props.check.status] || '#dadada';
        const borderColor = borderColorMap[this.props.check.status] || '#dfdfdf';
        const pie = css`
          width: 50%;
          height: 100%;
          transform-origin: 100% 50%;
          position: absolute;
          background: transparent;
          border: 5px solid ${borderColor};
        `
        const delay = this.props.check.prevRunSeconds;
        const duration = this.props.check.nextRunSeconds - delay;

        return (
            <div className={"checkSummary " + this.props.check.status} title={this.props.check.name} css={css`
                  background-color: ${color};
                `}>
                <div className="wrapper" css={css`
                  position: absolute;

                  width: 100%;
                  height: 100%;`}>
                    <div className="spinner" 
                        css={css`
                          ${pie};
                          border-radius: 100% 0 0 100% / 50% 0 0 50%;
                          z-index: 0;
                          border-right: none;
                          animation: ${this.state.animRotate} ${duration}s linear 1;
                          animation-delay: ${delay}s;
                          animation-fill-mode: forwards;
                    `}></div>
                    <div className="filler" 
                        css={css`
                          ${pie};
                          border-radius: 0 100% 100% 0 / 0 50% 50% 0;
                          left: 50%;
                          opacity: 0;
                          animation: ${this.state.animOpacity} ${duration}s steps(1, end) 1 reverse;
                          animation-delay: ${delay}s;
                          animation-fill-mode: forwards;
                          border-left: none;
                    `}></div>
                    <div className="mask"
                        css={css`
                          width: 50%;
                          height: 100%;
                          position: absolute;
                          border-radius: 100% 0 0 100% / 50% 0 0 50%;
                          background-color: ${color};
                          opacity: 1;
                          animation: ${this.state.animOpacity} ${duration}s steps(1, end) 1;
                          animation-delay: ${delay}s;
                          animation-fill-mode: forwards;
                    `}></div>
                </div>
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
                // var thing = this;
                // debugger;
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
                    {checks.map(check => <CheckSummary key={check.id} check={check} />)}
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
