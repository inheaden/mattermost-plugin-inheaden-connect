import React from 'react';

import {SVGS} from '../svgs';

export interface Props {}

/**
 *
 */
const Icon = ({}: Props) => {
    return (
        <i
            style={{
                display: 'grid',
                alignItems: 'center',
                justifyContent: 'center',
            }}
        >
            <svg
                style={{
                    fillRule: "evenodd",
                    clipRule: "evenodd",
                    strokeLinejoin: "round",
                    strokeMiterlimit: 2,
                }}
                viewBox="0 0 20 20"
            >
                <circle
                    cx={11.205}
                    cy={9.954}
                    r={7.302}
                    style={{
                        fill: "#fbfcfe",
                    }}
                    transform="matrix(1.36953 0 0 1.36953 -5.345 -3.633)"
                />
                <path
                    d="M3.985 10.001v-.002A3.214 3.214 0 0 1 7.184 6.8c.45 0 .895.095 1.305.278V2.861H2.217A2.181 2.181 0 0 0 .046 5.033v9.934c0 1.192.981 2.172 2.172 2.172h6.271v-4.217c-.41.183-.855.278-1.305.278a3.214 3.214 0 0 1-3.199-3.199Z"
                    style={{
                        fill: "#ff603f",
                        fillRule: "nonzero",
                    }}
                    transform="matrix(.78576 0 0 .79816 2.142 2.018)"
                />
                <path
                    d="M12.153 2.861H9.358v4.988a.438.438 0 0 1-.435.434.438.438 0 0 1-.273-.097 2.331 2.331 0 0 0-1.467-.519A2.343 2.343 0 0 0 4.85 9.999a2.344 2.344 0 0 0 2.333 2.332c.534 0 1.052-.183 1.467-.519a.432.432 0 0 1 .273-.097c.239 0 .435.197.435.435V17.139h2.794a2.184 2.184 0 0 0 2.173-2.172V5.033a2.182 2.182 0 0 0-2.172-2.172ZM18.726 5.201l-3.022 1.368a.87.87 0 0 0-.51.791v5.28c0 .341.2.651.511.792l3.021 1.367a.871.871 0 0 0 1.228-.792V5.992a.873.873 0 0 0-.869-.869.872.872 0 0 0-.359.078Z"
                    style={{
                        fill: "#ff603f",
                        fillRule: "nonzero",
                    }}
                    transform="matrix(.78576 0 0 .79816 2.142 2.018)"
                />
            </svg>
        </i>
    );
};

export default Icon;
