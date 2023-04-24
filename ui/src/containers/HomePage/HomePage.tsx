import React, { ReactElement } from 'react';

import PageSection from 'components/PageSection';
import UserServiceAccountToken from './UserServiceAccountToken';

export default function HomePage(): ReactElement {
  return (
    <PageSection header="Authenticating">
      <div className="md:w-1/2 mx-2">
        <UserServiceAccountToken />
      </div>
    </PageSection>
  );
}
