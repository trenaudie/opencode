
Here is a good representation of how to use SVGs.
The SVGs must be imported as a raw string from the public dir, using the '?raw' flag
Always use the following format import logo from '/public/<svgname>.svg?raw';
eg. import logo from '/public/logo.svg?raw';
CRITICAL: You must select the Path children of the SVG and fill them with a color. Use a FOR loop for this. See below.
import {
  Rect,
  Node,
  makeScene2D, SVG,Path
} from '@motion-canvas/2d';
import {
  createSignal,
  createRef,
  all,
  waitFor,
  linear,
  Vector2
} from '@motion-canvas/core';
import logo from '/public/logo.svg?raw';
console.log(`found logo ${logo}`)
const numRows = 5;
const baseWidth = 560;
const widthStep = 80;
const boxHeight = 48;
const gap = 8;

export default makeScene2D(function* (view) {
  // Set background color
  view.fill('#000000');

  // Refs for rectangles and signals for opacity
  const rectRefs = Array.from({ length: numRows }, () => createRef<Rect>());
  const rectSignals = Array.from({ length: numRows }, () => createSignal(0));
  const logoRef = createRef<SVG>();
  // Main container Node to stack rectangles
  const container = createRef<Node>();
  const gap_between_last_rectangle_and_svg_y = 250;
  view.add(
    <>
  <Node ref={container}>
      {Array.from({ length: numRows }, (_, i) => {
        const width = baseWidth - i * widthStep;
        const y = () => (i * (boxHeight + gap)) - (numRows * (boxHeight + gap) / 2) + (boxHeight / 2);

        return (
          <Rect
            ref={rectRefs[i]}
            width={width}
            height={boxHeight}
            fill={'#6cf1c2'}
            y={y}
            opacity={rectSignals[i]} // opacity controlled by signal
          />
        );
      })}
    </Node>
      <SVG ref={logoRef} svg={logo.replace("@color", "#f2ff48")} size = {300}  position= {() => {
        let svg_world_to_parent = logoRef().worldToParent();
      console.log(`the last child of the container is ${container().children().at(-1).absolutePosition()}`);
            console.log(`which, converting to the Parent of the SVG means: ${svg_world_to_parent.transformPoint(container().children().at(-1).absolutePosition())}`);
            let gap_between_last_rectangle_and_svg_vector =  new Vector2(0, gap_between_last_rectangle_and_svg_y);
            console.log(`gap is ${gap_between_last_rectangle_and_svg_vector}`)
    return container().children().at(-1).absolutePosition().transformAsPoint(svg_world_to_parent).add(gap_between_last_rectangle_and_svg_vector)
    }}/>
  </>
    
  );
  // Animate SVG elements
  const svgElements = [logoRef()]
  const animations = [];
  for (const svgElement of svgElements) {
    animations.push(svgElement.scale(1, .5));
    animations.push(svgElement.opacity(1, .5));
    for (const child of svgElement.children()[0].children()) {
      if (child instanceof Path) {
        yield* child.fill('white',1);
      }
    }
  }
  // Fade in rectangles one by one
  for (let i = 0; i < numRows; i++) {
    yield* all(
      rectSignals[i](1, 0.5, linear), // Fade to opaque
      waitFor(0.2) // Wait before next rectangle
    );
  }
});


This is also a good representation of how you can position points one to another, you must always give the x={} and y={} position relative to the PARENT element. So, for example, if you want object A to have a position that depends on object B that is not the parent, one option is:
1. to fetch the ABSOLUTE position of object B. eg. objectBAbsPos = objectBRef().absolutePosition()
2. convert it to the local coordinates of the PARENT of object A. eg. objectBAbsPos_inParentcoords = objectBAbsPos.transformAsPoint(objectARef().worldToParent())
3. add the offset that you want to create between the two objects. eg. objectAPos = objectBAbsPos_inParentcoords.add(new Vector2([0,200]))
`