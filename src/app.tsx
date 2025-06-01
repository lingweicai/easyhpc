/*
 * Part of HPC resources management project of easyhpc.cn developed published on 
 * https://github.com/lingweicai/easyhpc
 */

import React from 'react';
import cockpit from 'cockpit';
import { Page, PageSection } from '@patternfly/react-core';
import SinfoProvider from './SinfoContext';
import Dashbar from './Dashbar';
import Dashboard from './Dashboard';
import ExecutePanel from './execute/ExecutePannle';

const _ = cockpit.gettext;

export const Application = () => {
    return (
        <Page className='no-masthead-sidebar'>
            <PageSection
                style={{ backgroundColor: '#f2f2f2' }}
                aria-label="_(easyhpc section)"
                padding={{ default: 'noPadding' }}
                hasOverflowScroll
            >
                <SinfoProvider>
                    <Dashbar />
                    <Dashboard />
                    <ExecutePanel />
                </SinfoProvider>
            </PageSection>
        </Page>
    );
};
