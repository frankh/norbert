import React from 'react';
import { Query } from "react-apollo";
import { List, Icon } from "antd";
import { formatRelative } from 'date-fns'
import './App.css';

import { GET_CHECK } from "../queries";


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

class CheckResultsList extends React.PureComponent {
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


class CheckPage extends React.PureComponent {
    render() {
        return (
            <Query query={GET_CHECK} variables={{ checkId: this.props.match.params.id }}>
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
                        <div className="checkResultsList">
                            <List
                                itemLayout="horizontal"
                                dataSource={data.getCheck.results}
                                renderItem={result => (
                                    <CheckResultsList key={result.id} result={result} />
                                )}
                            />
                        </div>
                    )
                }}
            </Query>
        )
    }
}

export default CheckPage;
