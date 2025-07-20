{
  "title": "string (scene short name/label)",
  "description": "string (concise, what is shown/animated, key visual/logic concept)",

  "sceneMetadata": {
    "background": "string|null",                                // Used if explicitly set, else omit/null
    "canvasDefaults": "object|null",                            // As present (empty/default if not specified)
    "viewport": "string|null",                                  // Only if present in scene
    "canvasResolution": "object|null",                          // Only if present
    "other": "object|null"                                      // Misc, if needed (from scene objects)
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
      "ref": "string|null",                             // e.g. "colA", "circleRefs[0]", or null
      "parent": "string|null",                          // For top-level, value is "view", or another name in this list
      "properties": {
        // Only those present across all examples, Partial, as used
        "x": "number|string|null",
        "y": "number|string|null",
        "position": "[number,number]|string|null",
        "size": "number|[number,number]|null",
        "width": "number|string|null",
        "height": "number|string|null",
        "fill": "string|null",
        "stroke": "string|null",
        "lineWidth": "number|null",
        "start": "number|null",
        "end": "number|null",
        "scale": "number|null",
        "radius": "number|null",
        "grow": "number|null",
        "gap": "number|null",
        "padding": "number|null",
        "margin": "number|null",
        "alignItems": "string|null",
        "justifyContent": "string|null",
        "direction": "string|null",
        "layout": "boolean|null",
        "fontWeight": "number|string|null",
        "fontSize": "number|null",
        "fontFamily": "string|null",
        "opacity": "number|null",
        "text": "string|array|null",
        "tex": "string|array|null",
        "code": "string|null",
        "data": "string|null",
        "arrowSize": "number|null",
        "endArrow": "boolean|null",
        "lineCap": "string|null",
        "points": "array|null",
        "others": "object|null"
      },
      "anchor": "string|null",
      "staticStyling": "object|string|null",
      "children": "array|null",    // (Repeat objectInitialization format, supports deep nesting)
      "count": "number|null"       // If this object is a group/array

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
      "target": "string|array|null",     // ref name(s) or selector(s) of node(s)
      "property": "string|array|null",   // e.g. "fill", ["fill","stroke"], "opacity"

      // Value assignment
      "from": "number|string|array|object|null", 
                                        // optional start value (for clarity)
      "to": "number|string|array|object|null",   
                                        // required end value

      // Aliases (optional, mirror animationFlow)
      "startValue": "number|string|array|object|null",
      "endValue":   "number|string|array|object|null",

      // Documentation
      "action": "string|null",           // human summary, e.g. "color all paths yellow"
      "details": "object|string|null",   // any extra props or context
      "notes": "string|null"             // implementer hints
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
    "target": "string|array|null",       // Reference/name(s) of animated object(s)
    "property": "string|array|null",     // 'x', 'fill', ['x','y'], 'fillColor', etc

    // Value transitions
    "from": "number|string|array|object|null",   // Start value(s)
    "to": "number|string|array|object|null",     // End value(s)
    "startValue": "number|string|array|object|null", // Alias for compatibility
    "endValue": "number|string|array|object|null",

    // Duration / time
    "duration": "number|null",            // In seconds (support float/duration/named)
    "timingFunction": "string|null",      // Easing, 'easeInOutCubic' etc. Also can specify for springs.
    "easing": "string|null",              // Alias for compatibility
    "springType": "string|null",          // For spring animation ('PlopSpring', 'SmoothSpring', etc)

    // Tween/step logic
    "mapping": "string|null",             // e.g., 'map(start, end, easeFn(value))'
    "stepLogic": "string|object|null",    // Human summary, or actual mapped MC code for value updating
    "animationFn": "tween|spring|waitFor|edit|all|map|custom|null", // Explicit function if significant
    "call": "string|null",                // If 'callable', e.g., how the tween or spring is used

    // Control/concurrency
    "concurrency": "single|parallel|sequential|staggered|grouped|none|null", // Execution relation
    "runs": "description of sequence or concurrency",
    "order": "number|string|null",

    // Functional or advanced options
    "tweenCallback": "string|object|null",   // If a callback is used for custom value mapping in tween
    "easingCurve": "string|array|null",      // e.g. [custom array] if applicable
    "delays": "array|number|null",           // For staggered/concurrent tweens
    "labelsOrComments": "array|string|null", // Timeline comments, labels, debug highlights, as present

    // Waits and pauses
    "wait": "number|null",                   // Pause after/before, in seconds
    "waitType": "waitFor|yield|custom|null", // To clarify wait style if ambiguous

    // Additional
    "action": "string|null",                 // Human description ("fade out", "move left", etc)
    "details": "object|string|null",         // Any other MC-specific tweak, mapping, additional callback info.
    "notes": "string|null"                   // Freeform for implementer ("spring reverses direction", etc)
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
        "strategy": "string|null",
        "avgReturn": "string|null",
        "volatility": "string|null",
        "weekendHandling": "string|null"
      },
      "length": "number|null",
      "mapping": {
        "X axis (time)": "string|null",
        "Y axis (value)": "string|null"
      }
    }
    // ...add logic notes as present
  ],

  "enhancements": [
    // Optional, only if present in scene: e.g. labels, tracker annotations, timeline comments
    {
      "type": "label|tracker|comment|annotation",
      "location": "string|null",
      "text": "string|array|null",
      "logic": "string|array|null"
    }
    // ...add more categories as needed
  ],

  "OptionalEnhancements": {
    "loops": "false|object|string|null",
    "conditionals": "false|object|string|null",
    "eventTriggers": "false|object|string|null",
    "interactivity": "false|object|string|null",
    "timelineCommentsOrLabels": "false|array|object|null",
    "timelineAnnotations": "array|object|null",
    "consoleLogging": "array|object|null",
    "sceneMetadata": "object|null"
  },

  "notes": [
    // Freeform, add any additional clarification needed for recreation
    "string"
  ]
}