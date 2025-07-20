{
  "title": "string (scene short name/label)",
  "description": "string (concise, what is shown/animated, key visual/logic concept)",

  "sceneMetadata": {
    "background": "string",                                // Used if explicitly set, else omit/null
    "canvasDefaults": "object",                            // As present (empty/default if not specified)
    "viewport": "string",                                  // Only if present in scene
    "canvasResolution": "object",                          // Only if present
    "other": "object"                                      // Misc, if needed (from scene objects)
  },

  "imports": [
    {
      "type": "component|function|utility|reactive|utility/internal",
      "name": "string|array",
      "from": "string (import source module)"
    }
    {
      "type": "asset/svg-raw",
      "name": "string|array",             // e.g. "logo" or ["iconA","iconB"]
      "from": "string (must end in '?raw')" // e.g. "/public/logo.svg?raw"
    }
  ],

  "objectInitialization": [
    // All nodes, layouts, blocks, visuals, in order (may be nested with 'children')
    {
      "name": "string (optional, for referencing/hierarchy clarity)",
      "type": "component/class (Circle, Layout, CodeBlock, Rect, Node, Path, Grid, Line, Txt, Latex, etc.)",
      "ref": "string",                             // e.g. "colA", "circleRefs[0]", or null
      "parent": "string",                          // For top-level, value is "view", or another name in this list
      "properties": {
        "x": "number|string",
        "y": "number|string",
        "position": "[number,number]|string",
        "size": "number|[number,number]", // must use for Circle, for SVG
        "width": "number|string", //good for Rect, Layout 
        "height": "number|string",//good for Rect, Layout
        "fill": "string",
        "stroke": "string",
        "lineWidth": "number",
        "start": "number",
        "end": "number",
        "scale": "number",
        "grow": "number",
        "gap": "number",
        "padding": "number",
        "margin": "number",
        "alignItems": "string",
        "justifyContent": "string",
        "direction": "string",
        "layout": "boolean",
        "fontWeight": "number|string",
        "fontSize": "number",
        "fontFamily": "string",
        "opacity": "number",
        "text": "string|array",
        "tex": "string|array",
        "code": "string",
        "data": "string",
        "arrowSize": "number",
        "endArrow": "boolean",
        "lineCap": "string",
        "points": "array",
        "others": "object"
      },
      "anchor": "string",
      "staticStyling": "object|string",
      "children": "array",    // (Repeat objectInitialization format, supports deep nesting)
      "count": "number"       // If this object is a group/array

    }
    // ...more as needed by scene
  ],

  "refs": {
    // (Optional map) Keyed by ref name or variable, value is node/component type or usage
    "refName": "type/description"
  },

  "relativePositioning": {
    "hierarchy": [
      // Explicit or compact textual/array tree, top-down
      // E.g. ["view", ["Layout", ["Rect(colA)", ...]]]
      // or, dictionaries for clarity:
      {
        "parent": "string",
        "children": ["string", ...]
      }
      // ...nested/composed as needed for clarity
    ],
    "layoutLogic": [
      "string (explain flex/fill/grow/column rules, per layout, as in your grid/flex scenes)",
      "string (note, e.g., 'colA and colB grow proportionally, rowA vertical in column, etc.')"
    ],
    "alignment": [
      "string (describe horizontal/vertical alignment, default behaviors, notes)"
    ],
    "placement": [
      // Provide additional location/mapping info, especially for animated or axis-projected objects
      "string (e.g., 'Tracker Y matches moving S&P point', 'Label at (x,y), ...')"
    ],
    "referenceSummary": [
      // If used, summary for named refs/variables, to correlate with code
      "string"
    ]
  },
  "contentManipulation": [
    {
      "step": "number|string",           // order or phase (e.g. 1, "init", etc.)

      "type": "string",                  // e.g. "set" | "edit" | "custom"
                                        // analogous to animationFlow.type but no timing
      // Targets
      "target": "string|array",     // ref name(s) or selector(s) of node(s)
      "property": "string|array",   // e.g. "fill", ["fill","stroke"], "opacity"

      // Value assignment
      "from": "number|string|array|object", 
                                        // optional start value (for clarity)
      "to": "number|string|array|object",   
                                        // required end value

      // Aliases (optional, mirror animationFlow)
      "startValue": "number|string|array|object",
      "endValue":   "number|string|array|object",

      // Documentation
      "action": "string",           // human summary, e.g. "color all paths yellow"
      "details": "object|string",   // any extra props or context
      "notes": "string"             // implementer hints
    }
    // …add one dict per static change you need…
  ]

  }

  "animationFlow": [
    // Each animation step or phase as an object (order in array = time order in scene)
  {
    "step": "number|string",             // Sequence order or phase
    "type": "tween|spring|wait|custom",  // MC animation type

    // Targets
    "target": "string|array",       // Reference/name(s) of animated object(s)
    "property": "string|array",     // 'x', 'fill', ['x','y'], 'fillColor', etc

    // Value transitions
    "from": "number|string|array|object",   // Start value(s)
    "to": "number|string|array|object",     // End value(s)
    "startValue": "number|string|array|object", // Alias for compatibility
    "endValue": "number|string|array|object",

    // Duration / time
    "duration": "number",            // In seconds (support float/duration/named)
    "timingFunction": "string",      // Easing, 'easeInOutCubic' etc. Also can specify for springs.
    "easing": "string",              // Alias for compatibility
    "springType": "string",          // For spring animation ('PlopSpring', 'SmoothSpring', etc)

    // Tween/step logic
    "mapping": "string",             // e.g., 'map(start, end, easeFn(value))'
    "stepLogic": "string|object",    // Human summary, or actual mapped MC code for value updating
    "animationFn": "tween|spring|waitFor|edit|all|map|custom", // Explicit function if significant
    "call": "string",                // If 'callable', e.g., how the tween or spring is used

    // Control/concurrency
    "concurrency": "single|parallel|sequential|staggered|grouped|none", // Execution relation
    "runs": "description of sequence or concurrency",
    "order": "number|string",

    // Functional or advanced options
    "tweenCallback": "string|object",   // If a callback is used for custom value mapping in tween
    "easingCurve": "string|array",      // e.g. [custom array] if applicable
    "delays": "array|number",           // For staggered/concurrent tweens
    "labelsOrComments": "array|string", // Timeline comments, labels, debug highlights, as present

    // Waits and pauses
    "wait": "number",                   // Pause after/before, in seconds
    "waitType": "waitFor|yield|custom", // To clarify wait style if ambiguous

    // Additional
    "action": "string",                 // Human description ("fade out", "move left", etc)
    "details": "object|string",         // Any other MC-specific tweak, mapping, additional callback info.
    "notes": "string"                   // Freeform for implementer ("spring reverses direction", etc)
  }
  ],

  "signalsAndEffects": [
    // Only include if actually present in scene
    {
      "type": "Signal|Effect",
      "name": "string",
      "usage": "string (what/why/how used)"
    }
    // ...more as present
  ],

  "dataLogic": [
    // Only present in data-rich or chart/graph scenes
    {
      "syntheticData": {
        "strategy": "string",
        "avgReturn": "string",
        "volatility": "string",
        "weekendHandling": "string"
      },
      "length": "number",
      "mapping": {
        "X axis (time)": "string",
        "Y axis (value)": "string"
      }
    }
    // ...add logic notes as present
  ],

  "enhancements": [
    // Optional, only if present in scene: e.g. labels, tracker annotations, timeline comments
    {
      "type": "label|tracker|comment|annotation",
      "location": "string",
      "text": "string|array",
      "logic": "string|array"
    }
    // ...add more categories as needed
  ],

  "OptionalEnhancements": {
    "loops": "false|object|string",
    "conditionals": "false|object|string",
    "eventTriggers": "false|object|string",
    "interactivity": "false|object|string",
    "timelineCommentsOrLabels": "false|array|object",
    "timelineAnnotations": "array|object",
    "consoleLogging": "array|object",
    "sceneMetadata": "object"
  },

  "notes": [
    // Freeform, add any additional clarification needed for recreation
    "string"
  ]
}