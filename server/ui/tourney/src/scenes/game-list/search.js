import React, { Component } from 'react';

export default class Search extends Component {
    constructor(props) {
      super(props);
      this.handleFilterTextInputChange = this.handleFilterTextInputChange.bind(this);
    }
  
    handleFilterTextInputChange(e) {
      this.props.onFilterTextInput(e.target.value);
    }
  
    render() {
      return (
      <div className="input-group">
        <span className="input-group-addon">Filter</span>
        <input
              type="text"
              className="form-control"
              placeholder="Search..."
              value={this.props.filterText}
              onChange={this.handleFilterTextInputChange}
          />
      </div>
      );
    }
  }