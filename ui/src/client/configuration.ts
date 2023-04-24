import { Configuration, ConfigurationParameters } from '@stackrox/infra-auth-lib';

const parameters: ConfigurationParameters = {
  basePath: `${window.location.protocol}//${window.location.host}`,
};

export default new Configuration(parameters);
