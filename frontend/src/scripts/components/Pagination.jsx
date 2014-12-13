var React = require('react'),
    map = require('lodash-node/modern/collections/map');


var Pagination = React.createClass({
    propTypes: {
        // The current page information
        currPage:   React.PropTypes.number.isRequired,
        totalPages: React.PropTypes.number.isRequired,

        // Clicking to select a page.
        onClick: React.PropTypes.func,

        // Whether or not to center the paginator.
        center: React.PropTypes.bool,

        // Maximum pages to display.
        maxPages: React.PropTypes.number,
    },

    getDefaultProps: function() {
        return {
            center: true,
            maxPages: 10,
        };
    },

    handleClick: function(page, e) {
        console.log("Handling click on page: " + page);

        // Don't allow the event.
        this.preventClick(e);
    },

    preventClick: function(e) {
        e.preventDefault();
        e.stopPropagation();
    },

    render: function() {
        var havePrev = this.props.currPage > 1,
            haveNext = this.props.currPage < this.props.totalPages;

        // Shown if we don't have any pages
        var emptyPager = null;
        if( this.props.totalPages === 0 ) {
            emptyPager = (
              <li className="active">
                <a href="" onClick={this.preventClick}>1 <span className="sr-only">(current)</span></a>
              </li>
            );
        }

        // The maximum number of pages that we can display.
        var maxPages = this.props.totalPages > this.props.maxPages ?
                       this.props.maxPages :
                       this.props.totalPages;

        // Generate each item in the pager.
        var pagerItems = [];
        for( var i = 1; i <= maxPages; i++ ) {
            if( i === this.props.currPage ) {
                pagerItems.push(
                  <li className="active" key={i}>
                    <a href="" onClick={this.preventClick}>{i} <span className="sr-only">(current)</span></a>
                  </li>
                );
            } else {
                pagerItems.push(
                  <li key={i}>
                    <a href="" onClick={this.handleClick.bind(null, i)}>{i}</a>
                  </li>
                );
            }
        }

        // If there are more pages than we're displaying, then we display a disabled
        // "..." element at the end.
        if( this.props.totalPages > this.props.maxPages ) {
            pagerItems.push(
              <li className="disabled" key="dotdot">
                <a href="" onClick={this.preventClick}>...</a>
              </li>
            );
        }

        // Prev/next handler functions.
        var prevHandler = havePrev ?
                          this.handleClick.bind(null, this.props.currPage - 1) :
                          this.preventClick;
        var nextHandler = haveNext ?
                          this.handleClick.bind(null, this.props.currPage + 1) :
                          this.preventClick;

        return (
          <div className={this.props.center ? "text-center" : ""}>
            <nav>
              <ul className="pagination">
                {/* The "<<" arrow */}
                <li className={havePrev ? "" : "disabled"}>
                  <a href="" onClick={prevHandler}>&laquo;</a>
                </li>

                {/* If we don't have any documents, just show the current page, disabled */}
                {emptyPager}

                {/* The actual pager items */}
                {pagerItems}

                {/* The ">>" arrow */}
                <li className={haveNext ? "" : "disabled"}>
                  <a href="" onClick={nextHandler}>&raquo;</a>
                </li>
              </ul>
            </nav>
          </div>
        );
    },
});

module.exports = Pagination;
