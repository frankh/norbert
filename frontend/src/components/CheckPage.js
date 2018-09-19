import React from 'react';
import { Query } from "react-apollo";
import { List, Icon } from "antd";
import { formatRelative } from 'date-fns'

import { GET_CHECK, RESULT_SUBSCRIPTION } from "../queries";


class CheckResultIcon extends React.PureComponent {
    render() {
        const resultCode = this.props.resultCode;

        const iconMap = {
            "Success": "check-circle",
            "Failure": "exclamation-circle",
            "Error": "exclamation-circle",
        }
        const colorMap = {
            "Success": "#52c41a",
            "Failure": "#c41919",
            "Error": "#c41919",
        }

        return (
            <Icon type={iconMap[resultCode]} theme="twoTone" twoToneColor={colorMap[resultCode]} style={{fontSize: "32px"}}/>
        )
    }
}

class CheckResult extends React.PureComponent {
    render() {
        const result = this.props.result;
        return (
            <div className="checkResultItem">
                <List.Item>
                    <List.Item.Meta
                        avatar={<CheckResultIcon resultCode={result.resultCode} />}
                        title={result.resultCode}
                        description={formatRelative(new Date(result.startTime), new Date())}
                    />
                </List.Item>
            </div>
        )
    }
}


class LiveCheckPage extends React.PureComponent {
    onCheckResult = null;

    componentDidMount() {
        const { checkId, subscribe } = this.props;
        this.onCheckResult = subscribe({
            document: RESULT_SUBSCRIPTION,
            variables: { checkId: checkId },
            updateQuery: (prev, { subscriptionData: { data } }) => {
                const result = data.checkResultSub;

                var results = [result, ...prev.getCheck.results];

                return {
                    ...prev,
                    getCheck: {
                        ...prev.getCheck,
                        results: results,
                    }
                };
            },
        });
    }

    render() {
        return (
            <div className="checkResultsList">
                {this.props.children}
            </div>
        )
    }
}

class CheckPage extends React.PureComponent {
    render() {
        const checkId = this.props.match.params.id;
        return (
            <Query query={GET_CHECK} variables={{ checkId }}>
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
                        <LiveCheckPage checkId={checkId} subscribe={subscribeToMore}>
                            <List
                                itemLayout="horizontal"
                                dataSource={data.getCheck.results}
                                renderItem={result => (
                                    <CheckResult key={result.id} result={result} />
                                )}
                            />
                        </LiveCheckPage>
                    )
                }}
            </Query>
        )
    }
}

export default CheckPage;
