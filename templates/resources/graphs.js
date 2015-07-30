/* ********************************************************************************

	GRAPHS:

******************************************************************************** */

    google.setOnLoadCallback(drawCharts);

    var graph_width = 800;
	var graph_height = 100;

    function drawCharts() {
    	drawScoreChart();
    	drawTimeChart();
    	drawDepthChart();
    	drawTable();
    }

    var timeChart;
    function drawTimeChart() {
		var timechartDiv = document.getElementById('time_chart');
		var data = google.visualization.arrayToDataTable(times);
		var options = {
			curveType: 'function',
			width: graph_width,
			height: graph_height,
			chartArea: {width: '100%', height: '100%'},
			backgroundColor: 'black',
			legend: {position: 'none'},
			colors: [time_color[0],time_color[1]],
			vAxes: {
			  0: {title: 'Time', textStyle: {color: '999'}},
			},
			hAxis: {textPosition: 'in', textStyle: {color: '999'}, gridlines: {count: 0}}, 
			vAxis: {textPosition: 'in', gridlines: {color: time_grid}},
		};
		timeChart = new google.visualization.LineChart(timechartDiv);
		google.visualization.events.addListener(timeChart, 'select', function () {
	        var s = timeChart.getSelection();
	        if (s[0] != null) {
	        	var ply = ((s[0].row) * 2) + (s[0].column-1);
	        	selectPly(ply);
	    	}
		});
		timeChart.draw(data, options);
	}

	var depthChart;
	function drawDepthChart() {
		var depthchartDiv = document.getElementById('depth_chart');
		var data = google.visualization.arrayToDataTable(depths);
		var options = {
			curveType: 'function',
			width: graph_width,
			height: graph_height,
			chartArea: {width: '100%', height: '100%'},
			backgroundColor: 'black',
			legend: {position: 'none'},
			colors: [depth_color[0],depth_color[1]],
			hAxis: {textPosition: 'in', textStyle: {color: '#999'}, gridlines: {count: 0}},
			vAxes: {
		      0: {textPosition: 'none'},
		      1: {},
		    },
		    series: { 
              	0: {targetAxisIndex: 1},
              	1: {targetAxisIndex: 1}
            },
			vAxis: {textPosition: 'in', textStyle: {color: '#999'}, gridlines: {color: depth_grid}},
		};
		depthChart = new google.visualization.LineChart(depthchartDiv);
		google.visualization.events.addListener(depthChart, 'select', function () {
	        var s = depthChart.getSelection();
	        if (s[0] != null) {
	        	var ply = ((s[0].row) * 2) + (s[0].column-1);
	        	selectPly(ply);
	    	}
		});
		depthChart.draw(data, options);
	}

	var scoreChart;
	function drawScoreChart() {
		var data = google.visualization.arrayToDataTable(scores);
		var options = {
			width: graph_width,
			height: graph_height,
			chartArea: {width: '100%', height: '100%'},
			backgroundColor: 'black',
			legend: {position: 'none'},
			colors: [score_color[0],score_color[1]],
			hAxis: {textPosition: 'none', gridlines: {count: 0} },
			vAxis: {textPosition: 'in', textStyle: {color: '#999'}, gridlines: {color: score_grid}},
		};
		scoreChart = new google.visualization.ColumnChart(document.getElementById('score_chart'));
		google.visualization.events.addListener(scoreChart, 'select', function () {
	        var s = scoreChart.getSelection();
	        if (s[0] != null) {
	        	var ply = ((s[0].row) * 2) + (s[0].column-1);
	        	selectPly(ply);
	    	}
		});
		scoreChart.draw(data, options);
	}

	var table;
	function drawTable() {
        var data = google.visualization.arrayToDataTable(move_table_data);
        var cssClassNamesObj = {
        	headerRow: 'thClass',
        	tableRow: 'trClass',
        	oddTableRow: 'trOddClass',
        	selectedTableRow: 'trSelectedClass',
        	hoverTableRow: 'trHoverClass',
        	tableCell: 'cellClass',
        	headerCell: 'tdHeaderClass'
		};
		var options = {
			width: 444,
			height: 350,
			showRowNumber: false,
			allowHtml: true,
			sort: 'disable',
			cssClassNames: cssClassNamesObj
		};
        table = new google.visualization.Table(document.getElementById('move_table'));
        google.visualization.events.addListener(table, 'select', function () {
	        var s = table.getSelection();
	        if (s[0] != null) {
	        	selectPly(s[0].row);
	    	}
		});
        table.draw(data, options);
	}

	function selectPly(number) {
		setBoard(fens[number+1]);
		var data_row = parseInt(number / 2);
		var data_col = 2;
		if( number % 2 == 0 ) {
			data_col = 1;
		}
		depthChart.setSelection( [{row: data_row, column: data_col} ] );
		timeChart.setSelection( [{row: data_row, column: data_col} ] );
		table.setSelection( [ {row: number } ] );
	}