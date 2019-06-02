package mkolibaba.knightstour;

import javax.imageio.ImageIO;
import java.awt.Image;
import java.awt.Point;
import java.io.IOException;
import java.util.ArrayList;
import java.util.List;
import java.util.stream.Collectors;
import java.util.stream.Stream;

/**
 * @author Maksim Kolibaba
 * @since 02.06.2019
 */
public class Unit {
    public int x, y;
    public Image image;
    public List<Point> steppedTiles = new ArrayList<>();
    public int m, n;

    /**
     * Constructor.
     *
     * @param m     game.m value
     * @param n     game.n value
     * @param white if true then figure will be white colored, otherwise black
     */
    public Unit(int m, int n, boolean white) {
        this.m = m;
        this.n = n;
        try {
            image = ImageIO.read(Game.class.getResourceAsStream(String.format("/figure_%s.png", white ? "white" : "black")));
        } catch (IOException e) {
            e.printStackTrace();
        }
    }

    public void setUnitImageScale(int width, int height) {
        image = image.getScaledInstance(width, height, Image.SCALE_DEFAULT);
    }

    public void setPosition(int x, int y) {
        this.x = x;
        this.y = y;
    }

    public List<Point> getPossibleMovements() {
        return Stream.of(
                new Point(x - 1, y - 2),
                new Point(x + 1, y - 2),
                new Point(x - 2, y - 1),
                new Point(x + 2, y - 1),
                new Point(x - 2, y + 1),
                new Point(x + 2, y + 1),
                new Point(x - 1, y + 2),
                new Point(x + 1, y + 2)
        )
                .filter(p -> p.x >= 0 && p.x < m && p.y >= 0 && p.y < n)
                .filter(p -> !steppedTiles.contains(p))
                .collect(Collectors.toList());
    }
}
