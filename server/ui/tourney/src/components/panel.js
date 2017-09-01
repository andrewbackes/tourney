import React, { Component } from 'react';

class Panel extends Component {
    render() {
      return (
        <div className="panel-group">
          <div className={"panel panel-" + this.props.mode}>
            <div className="panel-heading">{this.props.title}</div>
            <div className="panel-body">
                {this.props.content}
            </div>
          </div>
        </div>
      );
    }
  }

export default Panel;