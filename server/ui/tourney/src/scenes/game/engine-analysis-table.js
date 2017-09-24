import React, { Component } from 'react';

export default class EngineAnalysisTable extends Component {
    render() {
      let rows = [];
      if (this.props.analysis) {
        this.props.analysis.forEach( (analysis, i) => {
          rows.unshift(<EngineAnalysisTableRow 
            key={i}
            raw={analysis}
          />);
        });
      }
      return (
        <div>
          <table className="table table-condensed table-fixed">
            <thead>
              <tr>
                <th className="col-xs-12">Raw</th>
                
              </tr>
            </thead>
            <tbody style={{ 'maxHeight' : '275px' }}>
              { rows }
            </tbody>
          </table>
        </div>
      );
    }
  }
  
  class EngineAnalysisTableRow extends Component {
    render() {
      return (
        <tr>
          <td className="col-xs-12">{this.props.raw}</td>
        </tr>
      );
    }
  }