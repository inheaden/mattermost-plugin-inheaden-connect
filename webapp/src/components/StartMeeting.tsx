import React from 'react';
import {connect} from 'react-redux';
import {bindActionCreators, Dispatch} from 'redux';

import {getCurrentChannelId} from 'mattermost-redux/selectors/entities/common';
import {makeStyleFromTheme} from 'mattermost-redux/utils/theme_utils';

import {startMeeting} from '../actions';
import {SVGS} from '../svgs';

const mapStateToProps = (state: any, ownProps: any) => ({
    ...ownProps,
    currentChannelId: getCurrentChannelId(state),
});

const dispatchProps = (dispatch: Dispatch) => {
    return {
        actions: bindActionCreators(
            {
                startMeeting,
            },
            dispatch,
        ),
    };
};

type mapProps = ReturnType<typeof mapStateToProps> &
ReturnType<typeof dispatchProps>;
export interface Props extends mapProps {
    creatorName: string;
    theme: any;
}

/**
 *
 */
const StartMeeting = ({
    actions,
    creatorName,
    currentChannelId,
    theme,
}: Props) => {
    const style = getStyle(theme);
    return (
        <div style={style.body}>
            <p>Has started a meeting:</p>
            <button
                className='btn btn-primary'
                style={style.button}
                onClick={() => actions.startMeeting(currentChannelId)}
            >
                <i
                    style={style.buttonIcon}
                    dangerouslySetInnerHTML={{
                        __html: SVGS.iconWhite,
                    }}
                />
                Join Meeting
            </button>
        </div>
    );
};

export default connect(mapStateToProps, dispatchProps)(StartMeeting);

const getStyle = makeStyleFromTheme((theme) => {
    return {
        body: {
            overflowX: 'auto',
            overflowY: 'hidden',
            paddingRight: '5px',
            width: '100%',
        },
        title: {
            fontWeight: '600',
        },
        button: {
            fontFamily: 'Jost, Open Sans',
            fontSize: '12px',
            fontWeight: 'bold',
            letterSpacing: '1px',
            lineHeight: '19px',
            marginTop: '12px',
            borderRadius: '4px',
            display: 'flex',
            alignItem: 'center',
            color: theme.buttonColor,
        },
        buttonIcon: {
            paddingRight: '8px',
            fill: theme.buttonColor,
            height: '19px',
        },
        summary: {
            fontFamily: 'Open Sans',
            fontSize: '14px',
            fontWeight: '600',
            lineHeight: '26px',
            margin: '0',
            padding: '14px 0 0 0',
        },
        summaryItem: {
            fontFamily: 'Open Sans',
            fontSize: '14px',
            lineHeight: '26px',
        },
    };
});
