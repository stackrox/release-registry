import React, { ReactElement } from 'react';
import { AxiosPromise } from 'axios';
import { Tooltip, TooltipOverlay } from '@stackrox/ui-components';

import useApiQuery from 'client/useApiQuery';
import configuration from 'client/configuration';
import { AlertCircle } from 'react-feather';
import { ClipLoader } from 'react-spinners';
import { ReleaseServiceApi, V1ReleaseServiceListResponse } from 'client/release';

const releaseService = new ReleaseServiceApi(configuration);
const fetchReleases = (): AxiosPromise<V1ReleaseServiceListResponse> =>
  releaseService.releaseServiceList();

export default function ListReleases(): ReactElement {
  const { loading, error, data } = useApiQuery(fetchReleases);

  if (loading) {
    return (
      <div className="inline-flex items-center">
        <ClipLoader size={16} color="currentColor" />
        <span className="ml-2">Loading releases ...</span>
      </div>
    );
  }

  if (error || !data?.releases) {
    return (
      <Tooltip content={<TooltipOverlay>{error?.message || 'Unknown error'}</TooltipOverlay>}>
        <div className="inline-flex items-center">
          <AlertCircle size={16} />
          <span className="ml-2">Unexpected error occurred while loading releases</span>
        </div>
      </Tooltip>
    );
  }

  const releases = data.releases.map((r) => (
    <div className="m-2 h-43">
      {r.tag} - {r.commit} - {r.creator}
    </div>
  ));

  return (
    <div>
      <h3 className="text-3xl mb-2">Releases</h3>
      <div className="flex items-center">
        <>{releases}</>
      </div>
    </div>
  );
}
