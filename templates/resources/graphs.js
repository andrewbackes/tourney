/* ********************************************************************************

	GRAPHS:

******************************************************************************** */

	var color_blue = '#009bd0';
	var color_lblue = '#cedef3';
	var color_dblue = '#003242';
	var color_purple = '#993CF3';
	var color_orange = '#e89800';
	var color_green = '#7CBC2D';
	var color_gray = '#808080';
	var color_dgray = '#333';
	
	var time_color = ['white', color_green];
	var depth_color = ['white', color_orange];
	var score_color = ['white', color_purple];

	var graph_width = 800;
	var graph_height = 100;
	var graph_background = "#2D3440";
	var graph_line_color = "#535863";
	var graph_axis_color = "#AAA";
	var graph_text_color = "#999";
	graph_number_color = "#999";

    google.setOnLoadCallback(drawCharts);
    
    function drawCharts() {
    	drawScoreChart();
    	drawTimeChart();
    	drawDepthChart();
    	drawTable();
    }
    
    function makeOptions(title, line_colors, axis_side) {
	    return {
			curveType: 'function',
			width: graph_width,
			height: graph_height,
			chartArea: {width: '100%', height: '100%'},
			backgroundColor: graph_background,
			legend: {position: 'none'},
			colors: [line_colors[0],line_colors[1]],
			axisTitlesPosition: 'in',
			vAxes: {
			  0: {title: title, textStyle: {color: graph_number_color}},
			},
			hAxis: {
				textPosition: 'in', 
				textStyle: {color: graph_text_color}, 
				gridlines: {count: 0}, 
				titleTextStyle: {color: graph_text_color},
			}, 
			vAxis: {
				textPosition: 'in', 
				gridlines: {color: graph_line_color},
				textStyle: {color: graph_text_color}, 
				baselineColor: graph_axis_color,
				titleTextStyle: {color: graph_text_color},
			},
			tooltip: {
				trigger: 'selection',
			},
		};
	}

    var timeChart;
    function drawTimeChart() {
		var timechartDiv = document.getElementById('time_chart');
		var data = google.visualization.arrayToDataTable(times);
		var options = makeOptions("Seconds Used", time_color, "left");
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
		/*
		var options = {
			curveType: 'function',
			width: graph_width,
			height: graph_height,
			chartArea: {width: '100%', height: '100%'},
			backgroundColor: graph_background,
			legend: {position: 'none'},
			colors: [depth_color[0],depth_color[1]],
			axisTitlesPosition: 'in',
			hAxis: {
				textPosition: 'in', 
				textStyle: {color: '#999'}, 
				gridlines: {count: 0}, 
				titleTextStyle: graph_text_color
			},
			vAxes: {
		      0: {title: 'Depth Reached', textStyle: {color: graph_number_color}},
		      1: {},
		    },
			vAxis: {textPosition: 'in', textStyle: {color: '#999'}, gridlines: {color: graph_line_color}, baselineColor: graph_axis_color},
			series: { 
              	0: {targetAxisIndex: 1},
              	1: {targetAxisIndex: 1}
            },
		};
		*/
		var options = makeOptions("Depth Reached", depth_color, "right");
		
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
		var options = makeOptions("Engine Evaluation", score_color, "left");
		options.vAxis.viewWindow = { max: 300, min: -300 };
		options.vAxis.gridlines.count = 7;
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
        	tableRow: 'trClass',
        	oddTableRow: 'trOddClass',
        	selectedTableRow: 'trSelectedClass',
        	hoverTableRow: 'trHoverClass',
        	tableCell: 'tdClass',
        	headerCell: 'thClass'
		};
		var options = {
			width: 437,
			height: 355,
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
		scoreChart.setSelection( [{row: data_row, column: data_col} ] );
		table.setSelection( [ {row: number } ] );
		
		var el = document.querySelector('#move_table > div > div:first-child');
		if (el) {
			var td = document.getElementsByClassName('tdClass');
			if ( td ) {
				var n = Math.max( 0 , number - 4 );
				el.scrollTop = td[0].offsetHeight * n ;
			}
		}
		
	}