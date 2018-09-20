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
