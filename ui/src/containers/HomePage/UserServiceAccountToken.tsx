import React, { ReactElement } from 'react';
import { AlertCircle, Clipboard } from 'react-feather';
import { ClipLoader } from 'react-spinners';
import { AxiosPromise } from 'axios';
import { useClipboard } from 'use-clipboard-copy';
import { Tooltip } from '@patternfly/react-core';
import { V1TokenResponse, UserServiceApi } from '@stackrox/infra-auth-lib';

import useApiQuery from 'client/useApiQuery';
import configuration from 'client/configuration';

const userService = new UserServiceApi(configuration);

const fetchToken = (): AxiosPromise<V1TokenResponse> => userService.userServiceToken({});

export default function UserServiceAccountToken(): ReactElement {
  const { loading, error, data } = useApiQuery(fetchToken);
  const clipboard = useClipboard({
    copiedTimeout: 800, // duration in milliseconds to show 'successfully copied' feedback
  });

  if (loading) {
    return (
      <div className="inline-flex items-center">
        <ClipLoader size={16} color="currentColor" />
        <span className="ml-2">Loading service account token...</span>
      </div>
    );
  }

  if (error || !data?.Token) {
    return (
      <Tooltip content={<div>{error?.message || 'Unknown error'}</div>}>
        <div className="inline-flex items-center">
          <AlertCircle size={16} />
          <span className="ml-2">
            Unexpected error occurred while loading service account token
          </span>
        </div>
      </Tooltip>
    );
  }

  return (
    <div>
      <h3 className="text-3xl mb-2">Authentication Token</h3>
      <div className="flex items-center">
        <p className="text-xl">Copy the following token for Bearer Authentication:</p>
        <button
          type="button"
          aria-label="Copy to clipboard"
          onClick={clipboard.copy}
          className="ml-2"
        >
          <Tooltip content={<div>Copy to clipboard</div>}>
            <div className="flex items-center">
              <Clipboard size={16} />
              {clipboard.copied && <span className="ml-2 text-success-700">Copied!</span>}
            </div>
          </Tooltip>
        </button>
      </div>
      <textarea
        rows={6}
        value={data.Token}
        className="mt-4 w-full bg-base-100 p-1 rounded border-2 border-base-300 hover:border-base-400 font-600 leading-normal outline-none"
        readOnly
        ref={clipboard.target}
      />
    </div>
  );
}
