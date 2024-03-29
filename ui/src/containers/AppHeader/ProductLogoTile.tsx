import React, { ReactElement } from 'react';
import { Link } from 'react-router-dom';
import { Flex, FlexItem } from '@patternfly/react-core';

import RHACSLogo from 'components/RHACSLogo';

export default function ProductLogoTile(): ReactElement {
  return (
    <Flex alignItems={{ default: 'alignItemsCenter' }}>
      <FlexItem>
        <Link to="/">
          <RHACSLogo />
        </Link>
      </FlexItem>
      <FlexItem>
        <span className="pf-u-font-family-redhatVF-heading-sans-serif pf-u-font-size-lg">
          Release Registry
        </span>
      </FlexItem>
    </Flex>
  );
}
