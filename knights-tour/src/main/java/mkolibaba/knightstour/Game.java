package mkolibaba.knightstour;

import javax.swing.JFrame;
import javax.swing.JOptionPane;
import java.awt.BorderLayout;
import java.awt.Canvas;
import java.awt.Color;
import java.awt.Dimension;
import java.awt.Graphics;
import java.awt.Point;
import java.awt.event.MouseAdapter;
import java.awt.event.MouseEvent;
import java.util.Random;

/**
 * @author Maksim Kolibaba
 * @since 01.06.2019
 */
public class Game extends Canvas {
    interface Colors {
        Color LIGHT_TILE_COLOR = new Color(247, 229, 211);
        Color DARK_TILE_COLOR = new Color(106, 63, 20);
        Color POSSIBLE_TILE = new Color(212, 230, 247);
        Color STEPPED_TILE = new Color(214, 214, 214);
    }

    private int m, n;

    private int width = 70;
    private int height = 70;

    private Unit unit;

    private boolean ended;

    public Game() {
        Random random = new Random();
        int mRange = 8;
        m = random.nextInt(mRange) + 5;
        int nRange = 8;
        n = random.nextInt(nRange) + 5;
    }

    public void start() {
        unit = new Unit(m, n, false);
        unit.setUnitImageScale((int) (width * 0.8), (int) (height * 0.8));
        Point start = new Point(0, 0);
        unit.setPosition(start.x, start.y);
        unit.steppedTiles.add(start);

        Dimension dimension = new Dimension(width * m, height * n);

        this.addMouseListener(new MouseAdapter() {
            @Override
            public void mousePressed(MouseEvent e) {
                moveTo(e.getX() / width, e.getY() / height);
            }
        });

        setMinimumSize(dimension);
        setMaximumSize(dimension);
        setPreferredSize(dimension);

        JFrame frame = new JFrame("Knight's Tour");
        frame.setDefaultCloseOperation(JFrame.EXIT_ON_CLOSE);
        frame.setLayout(new BorderLayout());
        frame.add(this, BorderLayout.CENTER);
        frame.setResizable(false);
        frame.pack();
        frame.setLocationRelativeTo(null);
        frame.setVisible(true);
    }

    private void moveTo(int x, int y) {
        Point point = new Point(x, y);
        if (ended || !unit.getPossibleMovements().contains(point)) {
            return;
        }

        unit.setPosition(x, y);
        unit.steppedTiles.add(point);
        repaint();

        if (unit.steppedTiles.size() == m * n) {
            JOptionPane.showMessageDialog(null, "You won! :)");
        } else if (unit.getPossibleMovements().isEmpty() && unit.steppedTiles.size() < m * n) {
            JOptionPane.showMessageDialog(null, "You lose: no possible movements available :(");
        }
    }

    @Override
    public void paint(Graphics g) {
        drawField(g);
        drawSteps(g);
        drawSteppedTo(g);
        drawUnit(g);
    }

    private void drawField(Graphics g) {
        for (int i = 0; i < m; i++) {
            for (int j = 0; j < n; j++) {
                drawRectangle(g, i, j, (i + j) % 2 == 0 ? Colors.DARK_TILE_COLOR : Colors.LIGHT_TILE_COLOR);
            }
        }
    }

    private void drawUnit(Graphics g) {
        g.drawImage(unit.image, unit.x * width + (int) (width * 0.1), unit.y * height + (int) (height * 0.1), null);
    }

    private void drawSteps(Graphics g) {
        unit.getPossibleMovements().forEach(point -> drawRectangle(g, point.x, point.y, Colors.POSSIBLE_TILE));
    }

    private void drawSteppedTo(Graphics g) {
        unit.steppedTiles.forEach(point -> drawRectangle(g, point.x, point.y, Colors.STEPPED_TILE));
    }

    private void drawRectangle(Graphics g, int x, int y, Color color) {
        g.setColor(color);
        g.fillRect(x * width, y * height, width, height);
    }
}

