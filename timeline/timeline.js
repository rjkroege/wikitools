/**
 * Basic file for defining the operation of the timeline.
 * Defines the bands and sets up the necessary stuff.
 */
Timeline.serverLocale = "en";
Timeline.clientLocale = "en";
Timeline.urlPrefix = "timeline/";

Timeline.GregorianDateLabeller.monthNames["en"] = [
    "Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];

Timeline.GregorianDateLabeller.dayNames["en"] = [
    "Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"];

Timeline.strings["en"] = {wikiLinkLabel:  "Read" };

var tl;
function onLoad() {
  var tl_el = document.getElementById("tl");
  var eventSource1 = new Timeline.DefaultEventSource();
  var theme1 = Timeline.ClassicTheme.create();

  // Set the Timeline's "width" automatically. Setting this on the first
  // band will affect all bands.
  theme1.autoWidth = true;
  theme1.timeline_start = new Date(Date.UTC(1990, 0, 1));
  theme1.timeline_stop  = new Date(Date.UTC(2020, 0, 1));
  
  var d = Timeline.DateTime.parseGregorianDateTime("2010");
  var bandInfos = [
    Timeline.createBandInfo({ // Make a month band
        width:          "70%", 
        intervalUnit:   Timeline.DateTime.WEEK,
        intervalPixels: 200,
        eventSource:    eventSource1,
        date:           d,
        theme:          theme1,
        layout:         'detailed'  // original, overview, detailed
     }),
     Timeline.createBandInfo({ // Make a year band
          width:          "30%", 
          intervalUnit:   Timeline.DateTime.YEAR, 
          intervalPixels: 100,
          eventSource:    eventSource1,
          date:           d,
          theme:          theme1,
          layout:         'overview'  // original, overview, detailed
     })];

  // Sync the bands...            
  bandInfos[1].syncWith = 0;
  bandInfos[1].highlight = true;                                                     

  // create the Timeline
  tl = Timeline.create(tl_el, bandInfos, Timeline.HORIZONTAL);
  
  // Set the base url for image, icon, etc. references.
  var url = '.';
  
  // Data is stored in the timeline_data variable.
  eventSource1.loadJSON(timeline_data, url);
  tl.layout();
}

var resizeTimerID = null;
function onResize() {
  if (resizeTimerID == null) {
    resizeTimerID = window.setTimeout(function() {
      resizeTimerID = null;
      tl.layout();
    }, 500);
  }
}
