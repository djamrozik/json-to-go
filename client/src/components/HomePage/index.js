import React, { useEffect, useState } from 'react';

import AceEditor from "react-ace";
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';
import { dark } from 'react-syntax-highlighter/dist/esm/styles/prism';

import Tooltip from '@material-ui/core/Tooltip';
import { CopyToClipboard } from 'react-copy-to-clipboard';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faClipboard } from '@fortawesome/free-solid-svg-icons'

import "ace-builds/webpack-resolver";
import "ace-builds/src-noconflict/mode-json";
import "ace-builds/src-noconflict/mode-golang";
import "ace-builds/src-noconflict/theme-monokai";

import { fetchConvertedJson, isJsonString } from './helpers';
import './index.css';

const currentDate = new Date();
const currentDayOfMonth = currentDate.getDate();
const currentDateStr = `${currentDate.getMonth()}-${currentDate.getDate()}-${currentDate.getFullYear()}`

const starterStruct = `type Generated struct {
    Info string \`json:"info"\`
    ExampleValues struct {
        DayOfMonth int \`json:"dayOfMonth"\`
        DateStr string \`json:"dateStr"\`
        SomeRandomFloat float64 \`json:"someRandomFloat"\`
        StrArray []string \`json:"strArray"\`
        IsWorking bool \`json:"isWorking"\`
    } \`json:"exampleValues"\`
}`

const starterJson = `{
  "info": "Add JSON here",
  "exampleValues": {
    "dayOfMonth": ${currentDayOfMonth},
    "dateStr": "${currentDateStr}",
    "someRandomFloat": ${Math.round(Math.random() * 10 * 100) / 100},
    "strArray": ["hello", "world"],
    "isWorking": true
  }
}`;

const ClipboardCopy = React.forwardRef((props, ref) => {
  return (
    <div {...props} ref={ref}>
      <CopyToClipboard text={props.text} onCopy={props.onCopy}>
        <FontAwesomeIcon icon={faClipboard} className="copy-clipboard" />
      </CopyToClipboard>
    </div>
  )
})

function HomePage() {
  const [errorMessage, setErrorMessage] = useState('');
  const [golangStruct, setGolangStruct] = useState(starterStruct);
  const [isRequestError, setIsRequestError] = useState(false);
  const [jsonText, setJsonText] = useState(starterJson);
  const [wasCopied, setWasCopied] = useState(false);

  useEffect(() => {
    if (wasCopied === false) {
      return;
    }
    setTimeout(() => {
      setWasCopied(false);
    }, 5000)
  }, [wasCopied])

  const onClipboardCopy = () => {
    setWasCopied(true);
  }

  const onJsonChange = newJson => {
    setJsonText(newJson);

    if (!isJsonString(newJson)) {
      return;
    }

    fetchConvertedJson(newJson, (golangStruct, errMessage) => {
      if (errMessage) {
        setErrorMessage(errMessage);
        setGolangStruct('');
        setIsRequestError(true);
        return;
      }
      setErrorMessage('');
      setIsRequestError(false);
      setGolangStruct(golangStruct);
    });
  }

  const isJsonValid = isJsonString(jsonText);

  const jsonValidInfoStyle = {
    'color': isJsonValid ? '#0ffc03' : 'red'
  }

  const editorStyle = {
    'width': '100%',
    'height': '100%'
  }

  const requestErrorStyle = {
    'color': 'red'
  }

  const tooltipTextStyle = {
    fontSize: '0.9rem',
    padding: '0.5rem 0.25rem',
    margin: '0'
  }

  return (
    <div className="home-page">
      <div className="home-page-header">
        <div className="header-cell">
          <div className="header-cell-left">
            JSON
          </div>
          <div className="header-cell-right" style={jsonValidInfoStyle}>
            {
              isJsonValid
                ? <span>JSON is Valid&nbsp;&#10004;</span>
                : <span>JSON is Not Valid</span>
            }
          </div>
        </div>
        <div className="header-cell">
          <div className="header-cell-left">
            Generated Golang Struct
          </div>
          <div className="header-cell-right" style={requestErrorStyle}>
            { isRequestError && 'Request Error' }
            { (!isRequestError && golangStruct) && (
              <div className="copy-clipboard-wrapper">
                {
                  wasCopied && <span className="copied-status-text">Copied&nbsp;&#10004;</span>
                }
                <Tooltip title={<p style={tooltipTextStyle}>Copy to Clipboard</p>}>
                  <ClipboardCopy text={golangStruct} onCopy={onClipboardCopy} />
                </Tooltip>
              </div>
            )}
          </div>
        </div>
      </div>
      <div className="home-page-content">
        <div className="home-page-content-cell">
          <AceEditor
            mode="json"
            theme="monokai"
            onChange={onJsonChange}
            value={jsonText}
            name="json-editor"
            style={editorStyle}
            fontSize={14}
            showPrintMargin={false}
            debounceChangePeriod={100}
          />
        </div>
        <div className="home-page-content-cell">
          <SyntaxHighlighter language={"go"} style={dark} className="resulting-struct">
            { isRequestError ? errorMessage : golangStruct }
          </SyntaxHighlighter>
        </div>
      </div>
      <div className="home-page-footer">
        json-to-golang @ {(new Date()).getFullYear()} /
        <a href="https://github.com/hyeres/json-to-golang">github</a> /
        <a href="https://github.com/hyeres/json-to-golang/issues/new">report an issue</a>
      </div>
    </div>
  )
}

export default HomePage
