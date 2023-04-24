import React, { ReactElement } from 'react';

import AppHeaderLayout from 'components/AppHeaderLayout';
import ProductLogoTile from './ProductLogoTile';
import UserInfo from './UserInfo';

export default function AppHeader(): ReactElement {
  const mainArea = <div className="flex h-full items-center ml-4" />;
  return <AppHeaderLayout logo={<ProductLogoTile />} main={mainArea} ending={<UserInfo />} />;
}
