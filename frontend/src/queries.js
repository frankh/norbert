import gql from "graphql-tag";

export const GET_SERVICES = gql`
    query GetServices {
        services {
          name
          url

          checks {
            id
            name
            severity
            status
            nextRunSeconds
            prevRunSeconds
          }
        }
    }
`;

export const SERVICE_SUBSCRIPTION = gql`
    subscription OnServiceUpdate($serviceName: String!) {
        serviceChanged(serviceName: $serviceName) {
          name
          url

          checks {
            id
            name
            severity
            status
            nextRunSeconds
            prevRunSeconds
          }
        }
    }
`;

export const GET_CHECK = gql`
    query GetCheck($checkId: String!) {
        getCheck(checkId: $checkId) {
          name

          results {
            id
            startTime
            endTime
            resultCode
            errorMsg
          }
        }
    }
`;

export const GET_RESULTS = gql`
    query GetResults($checkId: String!) {
        checkResults(checkId: $checkId) {
          id
          startTime
          endTime
          resultCode
          errorMsg
        }
    }
`;

export const RESULT_SUBSCRIPTION = gql`
    subscription CheckResultSub($checkId: String!) {
        checkResultSub(checkId: $checkId) {
          id
          startTime
          endTime
          resultCode
          errorMsg
        }
    }
`;
